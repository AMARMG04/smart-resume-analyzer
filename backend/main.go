package main

import (
	"bytes"
	"encoding/json"
	// "fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"strings"
	// "github.com/unidoc/unipdf/v3/common/license"
    // "github.com/unidoc/unipdf/v3/extractor"
    // "github.com/unidoc/unipdf/v3/model"
	"github.com/gin-gonic/gin"
	"github.com/ledongthuc/pdf"
)

func main() {
	router := gin.Default()

	// Allow CORS for localhost:3000
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "*")
		c.Next()
	})

	router.POST("/analyze", handleAnalyze)

	log.Println("Server running on http://localhost:8080")
	router.Run(":8080")
}

func handleAnalyze(c *gin.Context) {
	file, _, err := c.Request.FormFile("resume")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing resume file"})
		return
	}
	defer file.Close()

	jobDescription := c.PostForm("jobDescription")
	if jobDescription == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing job description"})
		return
	}

	resumeText, err := extractTextFromFile(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to extract resume text"})
		return
	}

	// log.Println("Full Resume Text:\n", resumeText)
	log.Println("Full Job Description:\n", jobDescription)

	const prompt = `You are an expert AI resume reviewer.

You will be given:
1. A candidate's resume (plain text)
2. A job description

Your task is to analyze them and return **detailed, structured feedback** using **clean markdown formatting**.

---

🧠 Your feedback must be insightful, realistic, and interview-focused — not generic. Every suggestion should help improve the candidate’s chances.

🎯 Especially focus on:
- Relevance to the job description
- Strength of keywords and phrasing
- Structure, clarity, formatting, and ATS-friendliness
- Measurable impact shown in projects and experience
- Areas that are outdated, vague, or too technical without purpose

---

📄 Return the result in this **strict markdown format**:

## 🔢 Match Percentage

85%

---

## ✅ Strengths

- **MERN Stack Expertise:** Candidate demonstrates deep understanding with full-stack projects like Byte Forge.
- **Next.js Experience:** Usage of modern SSR and frontend framework is aligned with current industry expectations.
- **Freelancing Experience:** Shows self-motivation and ability to deliver independently.
- Use clear bullet points, and bold the skill or keyword in each line.

---

## ⚠️ Suggestions to Improve

**Tailor Professional Summary:** The summary is generic. Incorporate keywords like "REST APIs," "responsive design," "database management," and the specific technologies mentioned in the job description. Frame it to emphasize the *value* you bring to *this* role.
**Quantify Byte Forge Impact:** While page load times and admin efficiency are mentioned, provide more specific numbers and context. For example, "Reduced page load times by up to 50%, resulting in a 15% increase in conversion rates." Quantify the "significant amount" with Razorpay.
**Expand Project Context:** Projects lack context. Briefly describe the *problem* each project solved and the specific design decisions made. Why Next.js? Why MongoDB? For IoT projects, mention data sources and potential challenges.
...

---

## 📝 Summary

Conclude with a short paragraph that summarizes the candidate's fit for the role, highlighting alignment, technical strengths, and one key area to work on.

---

❗️ Rules:
- Use markdown only — no HTML, no plain paragraphs.
- Use proper markdown headings like ##.
-Please format the "Suggestions to Improve" section as a properly numbered markdown list using "1", "2.", etc. Ensure each suggestion starts on a new line and avoid using bullet points ("-" or *) for this section. Do not include extra line breaks between list items.
- Use "-" (dash) for bullet points.
- Leave empty lines between all sections.
- Format output for clean rendering in React.

Return your analysis now.`
	

	fullPrompt := prompt + "\n\nResume:\n" + resumeText + "\n\nJob Description:\n" + jobDescription
	aiResponse, err := callOllama(fullPrompt)
	log.Println("AI Final Response:\n", aiResponse)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get AI response"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": aiResponse})
}

// func extractTextFromFile(file multipart.File) (string, error) {
//     var buf bytes.Buffer
//     io.Copy(&buf, file)

//     contentType := http.DetectContentType(buf.Bytes())
//     if !strings.HasPrefix(contentType, "application/pdf") {
//         // treat as plain text
//         return buf.String(), nil
//     }

//     // 🔐 Step 1: Load Unidoc License Key
//     err := license.SetMeteredKey("8658b8d10c3ff048539a94670bd23ffb421f70a9dd7e98f2350cb12323b144d5") // replace this!
//     if err != nil {
//         log.Println("License error:", err)
//         return "", err
//     }

//     // 🔍 Step 2: Load and extract text
//     reader := bytes.NewReader(buf.Bytes())
//     pdfReader, err := model.NewPdfReader(reader)
//     if err != nil {
//         return "", err
//     }

//     var textBuilder strings.Builder
//     numPages, err := pdfReader.GetNumPages()
//     if err != nil {
//         return "", err
//     }

//     for i := 1; i <= numPages; i++ {
//         page, err := pdfReader.GetPage(i)
//         if err != nil {
//             continue
//         }
//         ext, err := extractor.New(page)
//         if err != nil {
//             continue
//         }
//         text, err := ext.ExtractText()
//         if err == nil {
//             textBuilder.WriteString(text + "\n")
//         }
//     }

//     return textBuilder.String(), nil
// }

func extractTextFromFile(file multipart.File) (string, error) {
	var buf bytes.Buffer
	io.Copy(&buf, file)

	// Check if it's a PDF (simplified way)
	if strings.HasPrefix(http.DetectContentType(buf.Bytes()), "application/pdf") {
		reader, err := pdf.NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
		if err != nil {
			return "", err
		}
		var textBuilder strings.Builder
		pages := reader.NumPage()
		for i := 1; i <= pages; i++ {
			page := reader.Page(i)
			txt, _ := page.GetPlainText(nil)
			textBuilder.WriteString(txt)
		}
		return textBuilder.String(), nil
	}

	// Else treat as plain text
	return buf.String(), nil
}

// Call Ollama locally
func callOllama(prompt string) (string, error) {
	reqBody := map[string]string{
		"model":  "gemma3:12b",
		"prompt": prompt,
	}
	jsonData, _ := json.Marshal(reqBody)

	resp, err := http.Post("http://localhost:11434/api/generate", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var fullResponse strings.Builder
	decoder := json.NewDecoder(resp.Body)

	for decoder.More() {
		var chunk map[string]interface{}
		if err := decoder.Decode(&chunk); err != nil {
			break
		}
		if val, ok := chunk["response"].(string); ok {
			fullResponse.WriteString(val)
		}
	}

	return fullResponse.String(), nil
}