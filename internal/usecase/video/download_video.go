package video

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	db "github.com/dinhtatuanlinh/video/db/sqlc"
	internalError "github.com/dinhtatuanlinh/video/internal/error"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
	"time"
)

const (
	m3u8Type1 = "type1"
	m3u8Type2 = "type2"
	m3u8Type3 = "type3"
)

type DownloadVideoModel struct {
	Urls []UrlModel
}
type UrlModel struct {
	CategoryName string
	Name         string
	Url          string
}

func (u *UseCaseVideo) DownloadVideo(ctx context.Context, req *DownloadVideoModel) error {
	log.Info().Msg("DownloadVideo")
	if len(req.Urls) < 1 {
		err := errors.New("no url provided")
		log.Error().Err(err).Msg("Failed to create operator")
		return internalError.NewAppError("Failed to DownloadVideo", http.StatusBadRequest, codes.Internal, err)
	}
	outputDir := u.config.VideoPath
	_, err := os.Stat(outputDir)
	if os.IsNotExist(err) {
		os.MkdirAll(outputDir, 0755)
	}
	log.Info().Str("outputDir", outputDir).Msg("Output directory")
	for _, url := range req.Urls {
		var fileUrl string
		var folder string
		//fileUrl := "/downloads/"
		//folder := outputDir + "/"
		categoryName := url.CategoryName
		for {
			category, err := u.store.GetCategory(ctx, categoryName)
			if err != nil {
				if errors.Is(err, db.ErrRecordNotFound) {
					return internalError.NewAppError("operator_id not found", http.StatusBadRequest, codes.NotFound, err)
				}
				return internalError.NewAppError("internal error", http.StatusInternalServerError, codes.Internal, err)
			}
			categoryName = category.CategoryParentName
			fileUrl = category.VideoCategoryName + "/" + fileUrl
			folder = category.VideoCategoryName + "/" + folder
			if category.CategoryParentName == "" {
				break
			}
		}
		folder = outputDir + folder
		_, err := os.Stat(folder)
		if os.IsNotExist(err) {
			os.MkdirAll(folder, 0755)
		}
		fileUrl = "/downloads/" + fileUrl

		fmt.Println(">>>>>>>>>>>", fileUrl, folder)

		// Parse the URL
		//urlObj, err := neturl.Parse(url.Url)
		//if err != nil {
		//	return err
		//}

		// Get the base (last element) of the path
		//filename := path.Base(urlObj.Path)
		// Step 1: Select highest bandwidth variant
		variantURL, m3u8Type, err := getVariantPlaylist(url.Url)
		if err != nil {
			return err
		}
		log.Error().Err(err).Msgf("Selected highest bandwidth variant: %s", variantURL)

		// Step 2: Parse .m3u8 and download .ts files
		lines := []string{}
		if m3u8Type == m3u8Type2 {
			lines, _ = parseVariantAndDownload(variantURL, folder)
		} else if m3u8Type == m3u8Type1 {
			lines = downloadFileType1(variantURL, folder)
		}

		// Step 3: Create local playlist file
		name := strings.ReplaceAll(url.Name, " ", "_")
		playlistPath := path.Join(folder, name+".m3u8")
		err = os.WriteFile(playlistPath, []byte(strings.Join(lines, "\n")+"\n"), 0644)
		if err != nil {
			return err
		}
		//name := strings.Split(filename, ".")[0]
		request := db.CreateVideoParams{
			VideoCategoryName: url.CategoryName,
			Name:              name,
			Url:               fileUrl + name + ".m3u8",
		}
		_, err = u.store.CreateVideo(ctx, request)
		if err != nil {
			if db.ErrorCode(err) == db.UniqueViolation {
				log.Error().Err(err).Msg("Error unique violation")
				return internalError.NewAppError("Error unique violation", http.StatusBadRequest, codes.InvalidArgument, err)
			}
			log.Error().Err(err).Msg("Failed to create operator transaction")
			return err
		}
	}

	return nil
}

func getVariantPlaylist(masterURL string) (string, string, error) {
	resp, err := http.Get(masterURL)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)
	var maxBandwidth int
	var selectedPlaylist string
	var m3u8Type string
	var url string
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#EXT-X-STREAM-INF") {
			m3u8Type = m3u8Type2
			break
		} else if strings.HasSuffix(line, ".ts") {
			m3u8Type = m3u8Type3
			break
		} else if strings.HasPrefix(line, "http://") || strings.HasPrefix(line, "https://") {
			m3u8Type = m3u8Type1
			break
		}
	}
	if m3u8Type == m3u8Type2 {
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "#EXT-X-STREAM-INF") {
				// Extract BANDWIDTH
				var bandwidth int
				fmt.Sscanf(line, "#EXT-X-STREAM-INF:PROGRAM-ID=1,BANDWIDTH=%d", &bandwidth)

				// Read the next line for the playlist URL
				if scanner.Scan() {
					playlist := scanner.Text()
					if bandwidth > maxBandwidth {
						maxBandwidth = bandwidth
						selectedPlaylist = playlist
					}
				}
			}
		}
		if selectedPlaylist == "" {
			return "", "", fmt.Errorf("no stream found in master playlist")
		}
		url = resolveURL(masterURL, selectedPlaylist)
	} else if m3u8Type == m3u8Type1 {
		url = masterURL
	}

	return url, m3u8Type, nil
}

func resolveURL(baseURL, relative string) string {
	if strings.HasPrefix(relative, "http") {
		return relative
	}
	base := baseURL[:strings.LastIndex(baseURL, "/")+1]
	return base + relative
}

func downloadFile(url, filepath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}
func downloadFileType1(baseURL, destDir string) []string {
	// Get m3u8 content
	resp, err := http.Get(baseURL)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Prepare to read and store local lines
	scanner := bufio.NewScanner(resp.Body)
	var localLines []string
	segmentIndex := 0

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "#") {
			// Playlist info like #EXTINF
			localLines = append(localLines, line)
			continue
		}

		// This is a full URL to a .ts segment
		segmentURL := line
		now := time.Now()
		base := strings.ReplaceAll(now.Format("20060102150405.000"), ".", "")
		if !strings.HasSuffix(base, ".ts") {
			base += ".ts"
		}
		filename := fmt.Sprintf("%03d_%s", segmentIndex, base)
		filePath := path.Join(destDir, filename)

		fmt.Printf("📥 Downloading segment %d → %s\n", segmentIndex, filename)
		if err := downloadFile(segmentURL, filePath); err != nil {
			fmt.Printf("⚠️ Failed to download %s: %v\n", segmentURL, err)
			continue
		}

		localLines = append(localLines, filename)
		segmentIndex++
	}

	//// Save local playlist file
	//playlistPath := path.Join(destDir, "playlist.m3u8")
	//return os.WriteFile(playlistPath, []byte(strings.Join(localLines, "\n")+"\n"), 0644)
	return localLines
}

// Parse variant playlist and download segments concurrently
func parseVariantAndDownload(m3u8URL, destDir string) ([]string, []string) {
	resp, err := http.Get(m3u8URL)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var lines []string
	var tsFiles []string
	scanner := bufio.NewScanner(resp.Body)
	segmentIndex := 0

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasSuffix(line, ".ts") {
			tsURL := resolveURL(m3u8URL, line)

			// create unique filename like 000_original.ts
			now := time.Now()

			originalName := strings.ReplaceAll(now.Format("20060102150405.000"), ".", "") + ".ts"
			log.Info().Str("originalName", originalName).Msg("Downloading video")
			uniqueFile := fmt.Sprintf("%03d_%s", segmentIndex, originalName)
			savePath := path.Join(destDir, uniqueFile)

			tsFiles = append(tsFiles, uniqueFile)
			lines = append(lines, uniqueFile)

			// download directly without goroutines
			fmt.Printf("Downloading segment %d → %s\n", segmentIndex, uniqueFile)
			if err := downloadFile(tsURL, savePath); err != nil {
				fmt.Printf("Failed to download %s: %v\n", tsURL, err)
			}

			segmentIndex++
		} else {
			lines = append(lines, line)
		}
	}

	return lines, tsFiles
}
