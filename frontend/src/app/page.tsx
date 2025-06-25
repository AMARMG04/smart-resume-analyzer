'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';


export default function Home() {
  const [resumeFile, setResumeFile] = useState<File | null>(null);
  const [jobDescription, setJobDescription] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const router = useRouter();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!resumeFile || !jobDescription) {
      alert('Please upload a resume and enter a job description.');
      return;
    }

    const formData = new FormData();
    formData.append('resume', resumeFile);
    formData.append('jobDescription', jobDescription);

    setIsLoading(true);

    const res = await fetch('http://localhost:8080/analyze', {
      method: 'POST',
      body: formData,
    });

    const data = await res.json();
    console.log(data)
    localStorage.setItem('ai_result', data.result);
    setIsLoading(false);
    router.push('/result');
  };

  return (
    <main className="min-h-screen flex w-full justify-center items-center bg-gradient-to-br from-[#fafafa] to-[#e9e9e9] text-neutral-800 transition-colors duration-300">
      <div className="max-w-full mx-auto p-8">
        <div className="text-center mb-10">
          <h1 className="text-5xl font-semibold tracking-tight text-black">AI Resume Reviewer</h1>
          <p className="text-neutral-600 mt-3 text-lg">Instant feedback on your resume`s job match</p>
        </div>

        {isLoading ? (
          <div className="text-center py-20">
            <p className="text-xl font-medium animate-pulse">🔍 Analyzing your resume...</p>
          </div>
        ) : (
          <form
            onSubmit={handleSubmit}
            className="bg-black border border-neutral-200 rounded-2xl shadow-lg p-12 space-y-6 transition w-full"
          >
            <div>
              <label className="block text-sm font-medium text-white mb-2">Upload Resume</label>
              <input
                type="file"
                accept=".pdf,.txt"
                placeholder='hello'
                onChange={(e) => e.target.files && setResumeFile(e.target.files[0])}
                className="block w-full text-sm bg-neutral-700 text-white dark:text-white border border-neutral-300 dark:border-neutral-700 rounded-lg px-4 py-2 focus:outline-none focus:ring focus:ring-blue-400"
                required
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-white mb-2">Paste Job Description</label>
              <textarea
                rows={6}
                value={jobDescription}
                onChange={(e) => setJobDescription(e.target.value)}
                className="w-full resize-none text-sm bg-neutral-700 text-white placeholder:text-white rounded-lg px-4 py-3 focus:outline-none focus:ring focus:ring-blue-400"
                placeholder="Paste the job description here..."
                required
              />
            </div>

            <button
              type="submit"
              className="w-full bg-black dark:bg-white text-white dark:text-black font-semibold py-3 rounded-xl hover:opacity-90 transition"
            >
              Analyze Resume Fit
            </button>
          </form>
        )}
      </div>
    </main>
  );
}

// Create a new page file at `pages/result.tsx` for showing the result
