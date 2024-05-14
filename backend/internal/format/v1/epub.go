package format

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/htchan/BookSpider/internal/model"
)

func writeContainer(zipWriter *zip.Writer) error {
	containerFile, createErr := zipWriter.Create("META-INF/container.xml")
	if createErr != nil {
		return fmt.Errorf("create container file failed: %w", createErr)
	}

	_, writeErr := containerFile.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
	<container version="1.0" xmlns="urn:oasis:names:tc:opendocument:xmlns:container">
			<rootfiles>
					<rootfile full-path="OEBPS/content.opf" media-type="application/oebps-package+xml" />
			</rootfiles>
	</container>`))
	if writeErr != nil {
		return fmt.Errorf("write container file failed: %w", writeErr)
	}

	return nil
}

func writeMimeType(zipWriter *zip.Writer) error {
	mimeFile, createErr := zipWriter.Create("mimetype")
	if createErr != nil {
		return fmt.Errorf("create mimetype file failed: %w", createErr)
	}

	_, writeErr := mimeFile.Write([]byte(`application/epub+zip`))
	if writeErr != nil {
		return fmt.Errorf("write mimetype file failed: %w", writeErr)
	}

	return nil
}

func writeToc(zipWriter *zip.Writer, bk *model.Book, chapters model.Chapters) error {
	tocFile, createErr := zipWriter.Create("toc.ncx")
	if createErr != nil {
		return fmt.Errorf("create toc file failed: %w", createErr)
	}

	chaptersContent := ""
	for i, chapter := range chapters {
		chaptersContent += fmt.Sprintf(
			`<navPoint id="num_%d" playOrder="%d">
			<navLabel><text>%s</text></navLabel>
			<content src="OEBPS/chapters/chapter-%d.xhtml"/>
			</navPoint>`,
			i+1, i+1, chapter.Title, i+1,
		)
	}

	_, writeErr := tocFile.Write([]byte(fmt.Sprintf(
		`<?xml version="1.0" encoding="UTF-8"?>
		<ncx version="2005-1" xml:lang="zh" xmlns="http://www.daisy.org/z3986/2005/ncx/">
			<head>
				<meta name="dtb:depth" content="2"/> <!-- 1 or higher -->
				<meta name="dtb:totalPageCount" content="0"/> <!-- must be 0 -->
				<meta name="dtb:maxPageNumber" content="0"/> <!-- must be 0 -->
			</head>
		
			<docTitle>
				<text>%s</text>
			</docTitle>
		
			<navMap>
				%s
			</navMap>
		</ncx>`,
		bk.Title, chaptersContent,
	)))
	if writeErr != nil {
		return fmt.Errorf("write toc file failed: %w", writeErr)
	}

	return nil
}

func writeCover(zipWriter *zip.Writer, bk *model.Book) error {
	coverFile, createErr := zipWriter.Create("OEBPS/cover.xhtml")
	if createErr != nil {
		return fmt.Errorf("create cover file failed: %w", createErr)
	}

	_, writeErr := coverFile.Write([]byte(fmt.Sprintf(
		`<?xml version='1.0' encoding='utf-8'?>
		<html xmlns="http://www.w3.org/1999/xhtml">
		<body>
			<center>
				<h1>%s</h1>
				<hr/>
				<h2>%s</h2>
			</center>
		</body>
		</html>`,
		bk.Title, bk.Writer.Name,
	)))
	if writeErr != nil {
		return fmt.Errorf("write cover file failed: %w", writeErr)
	}

	return nil
}

func writeContent(zipWriter *zip.Writer, bk *model.Book, chapters model.Chapters) error {
	contentFile, createErr := zipWriter.Create("OEBPS/content.opf")
	if createErr != nil {
		return fmt.Errorf("create content file failed: %w", createErr)
	}
	manifestContent := ""
	spineContent := ""

	for i := range chapters {
		manifestContent += fmt.Sprintf(
			`<item id="chapter-%d" href="chapters/chapter-%d.xhtml" media-type="application/xhtml+xml" />`,
			i+1, i+1,
		)
		spineContent += fmt.Sprintf(`<itemref idref="chapter-%d" />`, i+1)
	}

	_, writeErr := contentFile.Write([]byte(fmt.Sprintf(
		`<?xml version="1.0" encoding="UTF-8"?>
		<package xmlns="http://www.idpf.org/2007/opf" version="3.0">
				<metadata xmlns:dc="http://purl.org/dc/elements/1.1/">
						<dc:title>%s</dc:title>
						<dc:creator xmlns:ns0="http://www.idpf.org/2007/opf" ns0:role="aut" ns0:file-as="%s">%s</dc:creator>
						<dc:language>zh</dc:language>
				</metadata>
				<manifest>
						%s
						<item id="toc" href="../toc.ncx" media-type="application/x-dtbncx+xml"/>
						<item id="cover-page" href="cover.xhtml" media-type="application/xhtml+xml" properties="calibre:title-page"/>
				</manifest>
				<spine toc="toc">
				<itemref idref="cover-page"/>
				%s
				</spine>
		</package>`,
		bk.Title, bk.Writer.Name, bk.Writer.Name, manifestContent, spineContent,
	)))
	if writeErr != nil {
		return fmt.Errorf("write content file failed: %w", writeErr)
	}

	return nil
}

func reformatChaptersForEpub(lines []string) []string {
	result := make([]string, 0)

	for _, line := range lines {
		data := strings.TrimSpace(line)
		if len(data) > 0 {
			result = append(result, fmt.Sprintf("<p>%s</p>", data))
		}
	}

	return result
}

func writeChapters(zipWriter *zip.Writer, chapters model.Chapters) error {
	for i, chapter := range chapters {
		chapterFile, createErr := zipWriter.Create(fmt.Sprintf("OEBPS/chapters/chapter-%d.xhtml", i+1))
		if createErr != nil {
			return fmt.Errorf("create chapter %d file failed: %w", i+1, createErr)
		}

		content := strings.Split(chapter.Content, "\n")
		content = reformatChaptersForEpub(content)

		_, writeErr := chapterFile.Write([]byte(fmt.Sprintf(
			`<?xml version="1.0" encoding="UTF-8"?>
			<!DOCTYPE html>
			<html xmlns="http://www.w3.org/1999/xhtml">
			<head>
					<title>%s</title>
			</head>
			<body>
			<h1>%s</h1>
			%s
			</body>
			</html>`,
			chapter.Title, chapter.Title, strings.Join(content, "\n"),
		)))
		if writeErr != nil {
			return fmt.Errorf("write chapter %d failed: %w", i+i, writeErr)
		}
	}

	return nil
}

func (serv *serviceImpl) WriteBookEpub(ctx context.Context, bk *model.Book, chapters model.Chapters, writer io.Writer) error {
	// prepare the zip file
	zipWriter := zip.NewWriter(writer)
	defer zipWriter.Close()

	writeContainer(zipWriter)
	writeMimeType(zipWriter)
	writeToc(zipWriter, bk, chapters)
	writeCover(zipWriter, bk)
	writeContent(zipWriter, bk, chapters)
	writeChapters(zipWriter, chapters)

	return nil
}
