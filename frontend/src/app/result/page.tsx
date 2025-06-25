'use client';

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import ReactMarkdown from 'react-markdown';

function decodeHTMLEntities(text: string) {
    const textarea = document.createElement('textarea');
    textarea.innerHTML = text;
    return textarea.value;
}

export default function ResultPage() {
    const [result, setResult] = useState<string | null>(null);
    const [strengths, setStrengths] = useState<string | null>(null);
    const [suggestions, setSuggestions] = useState<string[] | null>(null);
    const [summary, setSummary] = useState<string | null>(null);
    const [matchPercent, setMatchPercent] = useState<string | null>(null);
    const router = useRouter();

    useEffect(() => {
        const stored = localStorage.getItem('ai_result');
        if (!stored) return;

        const decoded = decodeHTMLEntities(stored);

        const matchPercentMatch = decoded.match(/## 🔢 Match Percentage\n\n(\d+%)/);
        const strengthsMatch = decoded.match(/## ✅ Strengths\n([\s\S]*?)---/);
        const suggestionsMatch = decoded.match(/## ⚠️ Suggestions to Improve\n([\s\S]*?)---/);
        const summaryMatch = decoded.match(/## 📝 Summary\n([\s\S]*)/);

        setMatchPercent(matchPercentMatch?.[1]?.trim() ?? null);
        setStrengths(strengthsMatch?.[1]?.trim() ?? null);

        if (suggestionsMatch?.[1]) {
            const raw = suggestionsMatch[1].trim();
            const suggestionList = raw
                .split(/\n(?=\w)/)
                .map(line => line.replace(/^\d+\.\s*/, '').trim()) // Remove leading "1. ", "2. ", etc.
                .filter(Boolean);
            setSuggestions(suggestionList);
        }

        setSummary(summaryMatch?.[1]?.trim() ?? null);
        setResult(decoded);
    }, []);

    return (
        <main className="min-h-screen bg-gradient-to-br from-[#fafafa] to-[#e9e9e9] text-neutral-800 transition-colors duration-300">
            <div className="max-w-3xl mx-auto p-8">
                <div className="text-center mb-10">
                    <h1 className="text-5xl font-bold tracking-tight text-black">AI Feedback</h1>
                    <p className="text-neutral-600 mt-3 text-lg">Here`s how your resume aligns with the job</p>
                    {matchPercent && (
                        <p className="text-4xl font-bold mt-4 text-green-600 dark:text-green-400">🎯 {matchPercent} Match</p>
                    )}
                </div>

                {result ? (
                    <>
                        {strengths && (
                            <section className="mb-10">
                                <h2 className="text-3xl font-semibold mb-4 text-green-600 dark:text-green-400">✅ Strengths</h2>
                                <div className="prose dark:prose-invert text-lg">
                                    <ReactMarkdown>{strengths}</ReactMarkdown>
                                </div>
                            </section>
                        )}

                        {suggestions && suggestions.length > 0 && (
                            <section className="mb-10">
                                <h2 className="text-3xl font-semibold mb-4 text-yellow-600 dark:text-yellow-400">⚠️ Suggestions to Improve</h2>
                                <ol className="list-decimal ml-6 space-y-4 text-lg leading-relaxed">
                                    {suggestions.map((point, idx) => (
                                        <li key={idx}><ReactMarkdown>{point}</ReactMarkdown></li>
                                    ))}
                                </ol>
                            </section>
                        )}

                        {summary && (
                            <section className="mb-10">
                                <h2 className="text-3xl font-semibold mb-4 text-blue-600 dark:text-blue-400">📝 Summary</h2>
                                <div className="prose dark:prose-invert text-lg">
                                    <ReactMarkdown>{summary}</ReactMarkdown>
                                </div>
                            </section>
                        )}

                        <div className="text-center mt-12">
                            <button
                                onClick={() => router.push('/')}
                                className="bg-black text-white font-semibold px-6 py-3 rounded-lg hover:opacity-90 transition"
                            >
                                🔁 Analyze Another Resume
                            </button>
                        </div>
                    </>
                ) : (
                    <p className="text-center text-neutral-500 dark:text-neutral-400">
                        No result found. Please go back and analyze your resume.
                    </p>
                )}
            </div>
        </main>
    );
}