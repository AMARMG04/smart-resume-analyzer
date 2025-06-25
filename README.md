# Smart Resume Analyzer 🧠📄

An AI-powered resume analysis tool that evaluates your resume against a job description and provides tailored feedback.

## 🔧 Tech Stack

- **Frontend:** Next.js, React, TailwindCSS
- **Backend:** Go (Golang)
- **AI Model:** Ollama (Google's Gemma3:12B) (via local)
- **Storage:** Local file parsing + HTML entity decoding

## 📁 Project Structure
smart-resume-analyzer/
├── frontend/   # Next.js frontend
├── backend/    # Go API server

## 🚀 How to Run

### 1. Frontend

```bash
cd frontend
npm install
npm run 
```

### 2. Backend


```bash
cd backend
go run main.go
```


## 📸 Screenshots
<img width="1491" alt="Screenshot 2025-06-25 at 9 25 53 PM" src="https://github.com/user-attachments/assets/3f37d9d1-4db9-4832-9325-bdcc36ffdddf" />
<img width="1491" alt="Screenshot 2025-06-25 at 9 26 09 PM" src="https://github.com/user-attachments/assets/abe00200-66a2-4151-b72c-9d28445087eb" />
<img width="1488" alt="Screenshot 2025-06-25 at 9 28 11 PM" src="https://github.com/user-attachments/assets/e5dfbc29-2df9-44b0-9767-db29ea20876a" />
<img width="1488" alt="Screenshot 2025-06-25 at 9 28 30 PM" src="https://github.com/user-attachments/assets/9299e1ea-8ad3-4e20-a30a-dbc1a7eb8211" />
<img width="1488" alt="Screenshot 2025-06-25 at 9 28 52 PM" src="https://github.com/user-attachments/assets/11957f38-aea7-4d24-a4c2-a9e209d8827c" />

## ✨ Features
	•	Upload resume and paste job description
	•	AI-generated match feedback (strengths, improvements, summary)
	•	Markdown rendering
	•	Mobile responsive UI

## 📦 Future Improvements
	•	Auth system for saving feedback
	•	Resume templates
	•	Job scraping integration
