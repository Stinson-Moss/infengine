package obsidian

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Stinson-Moss/infengine/db/postgres/db"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func ExportSourceDocument(doc db.Document, vaultPath string) error {
	safeTitle := sanitizeFilename(doc.Title)
	if safeTitle == "" {
		
		safeTitle = fmt.Sprintf("Untitled Document %d", doc.ID)
	}
	
	tCaser := cases.Title(language.Und)
	filename := fmt.Sprintf("%s.md", safeTitle)
	
	outputDir := filepath.Join(vaultPath, "01-Sources")
	filePath := filepath.Join(outputDir, filename)
	
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", outputDir, err)
	}

	var formattedDate string
	if doc.Created.Valid {
		formattedDate = doc.Created.Time.Format(time.RFC3339)
	} else {
		formattedDate = time.Now().Format(time.RFC3339) 
	}

	var md strings.Builder

	md.WriteString("---\n")
	md.WriteString(fmt.Sprintf("title: %q\n", tCaser.String(doc.Title)))
	md.WriteString(fmt.Sprintf("created_at: %s\n", formattedDate))
	md.WriteString("---\n\n")

	
	md.WriteString(fmt.Sprintf("# %s\n\n", tCaser.String(doc.Title)))

	
	if strings.TrimSpace(doc.Description) != "" {
		md.WriteString("## Description\n")
		
		escapedDesc := strings.ReplaceAll(doc.Description, "\n", "\n> ")
		md.WriteString(fmt.Sprintf("> [!abstract] Snippet\n> %s\n\n", escapedDesc))
	}

	
	md.WriteString("## Content\n")
	md.WriteString(doc.Content)
	md.WriteString("\n")

	
	err := os.WriteFile(filePath, []byte(md.String()), 0644)
	if err != nil {
		return fmt.Errorf("failed writing file %s: %w", filename, err)
	}

	return nil
}

func sanitizeFilename(name string) string {
	replacer := strings.NewReplacer(
		"/", "-", 
		"\\", "-", 
		":", " -", 
		"*", "", 
		"?", "", 
		"\"", "", 
		"<", "", 
		">", "", 
		"|", "",
		"[", "",
		"]", "",
		"#", "",
	)

	return strings.TrimSpace(replacer.Replace(name))
}