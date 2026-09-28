"use client";

import { useState } from "react";
import Link from "next/link";
import { Film } from "lucide-react";
import { api } from "@/lib/api/client";

export default function ForgotPasswordPage() {
  const [email, setEmail] = useState("");
  const [submitted, setSubmitted] = useState(false);
  // Development mode returns the token directly (no mailer).
  const [devToken, setDevToken] = useState<string | null>(null);
  const [isSubmitting, setIsSubmitting] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!email) return;
    setIsSubmitting(true);
    try {
      const res = await api.post<{ ok: boolean; reset_token?: string }>(
        "/v1/auth/forgot-password",
        { email }
      );
      setDevToken(res.reset_token ?? null);
      setSubmitted(true);
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="min-h-screen bg-studio-bg flex items-center justify-center p-8">
      <div className="w-full max-w-sm space-y-8">
        <div className="flex items-center gap-3 justify-center">
          <div className="w-10 h-10 rounded-md bg-studio-panel border border-studio-border flex items-center justify-center text-studio-accent">
            <Film className="w-5 h-5" />
          </div>
          <span className="font-bold text-xl tracking-tight text-white">DramaStudio</span>
        </div>

        <div className="bg-studio-card border border-studio-border rounded-md p-6 space-y-4">
          {submitted ? (
            <div className="space-y-4">
              <div>
                <h1 className="text-lg font-bold text-white">Check your email</h1>
                <p className="text-xs text-studio-muted mt-1">
                  If an account exists for {email}, a reset link is on its way.
                </p>
              </div>
              {devToken && (
                <div className="bg-studio-panel border border-studio-border rounded p-3 space-y-2">
                  <p className="text-[10px] text-amber-400 font-mono uppercase tracking-wider">
                    Development mode — token returned directly
                  </p>
                  <code className="block text-xs text-white break-all">{devToken}</code>
                  <Link
                    href={`/reset-password?token=${encodeURIComponent(devToken)}`}
                    className="inline-block text-xs text-studio-accent hover:underline"
                  >
                    Continue to reset password →
                  </Link>
                </div>
              )}
              <p className="text-xs text-studio-muted text-center">
                <Link href="/login" className="text-studio-accent hover:underline">
                  Back to sign in
                </Link>
              </p>
            </div>
          ) : (
            <form onSubmit={handleSubmit} className="space-y-4">
              <div>
                <h1 className="text-lg font-bold text-white">Reset your password</h1>
                <p className="text-xs text-studio-muted mt-1">
                  Enter your account email and we&apos;ll send a reset link.
                </p>
              </div>

              <div className="space-y-1">
                <label className="block text-xs font-semibold text-studio-muted">Email</label>
                <input
                  type="email"
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  placeholder="you@studio.com"
                  autoComplete="email"
                  className="w-full bg-studio-panel border border-studio-border rounded px-3 py-2 text-sm text-white placeholder:text-studio-muted/50 focus:outline-none focus:border-studio-accent"
                />
              </div>

              <button
                type="submit"
                disabled={isSubmitting || !email}
                className="w-full bg-studio-accent hover:bg-studio-accent-dark disabled:opacity-50 text-white font-medium text-sm py-2.5 rounded transition-colors"
              >
                {isSubmitting ? "Sending..." : "Send reset link"}
              </button>

              <p className="text-xs text-studio-muted text-center">
                Remembered it?{" "}
                <Link href="/login" className="text-studio-accent hover:underline">
                  Sign in
                </Link>
              </p>
            </form>
          )}
        </div>
      </div>
    </div>
  );
}
