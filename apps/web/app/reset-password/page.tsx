"use client";

import { Suspense, useState } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import Link from "next/link";
import { Film, Eye, EyeOff } from "lucide-react";
import { api } from "@/lib/api/client";

function ResetPasswordForm() {
  const router = useRouter();
  const params = useSearchParams();
  const [token, setToken] = useState(params.get("token") ?? "");
  const [password, setPassword] = useState("");
  const [showPassword, setShowPassword] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [isSubmitting, setIsSubmitting] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!token || password.length < 8) return;
    setIsSubmitting(true);
    setError(null);
    try {
      await api.post("/v1/auth/reset-password", { token, password });
      router.push("/login?reset=1");
    } catch (err) {
      setError(err instanceof Error ? err.message : "Reset failed.");
      setIsSubmitting(false);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="bg-studio-card border border-studio-border rounded-md p-6 space-y-4">
      <div>
        <h1 className="text-lg font-bold text-white">Set a new password</h1>
        <p className="text-xs text-studio-muted mt-1">
          Paste your reset token and choose a new password.
        </p>
      </div>

      {error && (
        <div className="bg-red-500/10 border border-red-500/30 text-red-400 rounded px-3 py-2 text-xs">
          {error}
        </div>
      )}

      <div className="space-y-1">
        <label className="block text-xs font-semibold text-studio-muted">Reset token</label>
        <input
          type="text"
          value={token}
          onChange={(e) => setToken(e.target.value)}
          placeholder="Paste the token from your reset email"
          autoComplete="off"
          className="w-full bg-studio-panel border border-studio-border rounded px-3 py-2 text-sm text-white font-mono placeholder:text-studio-muted/50 focus:outline-none focus:border-studio-accent"
        />
      </div>

      <div className="space-y-1">
        <label className="block text-xs font-semibold text-studio-muted">New password</label>
        <div className="relative">
          <input
            type={showPassword ? "text" : "password"}
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            placeholder="Minimum 8 characters"
            autoComplete="new-password"
            className="w-full bg-studio-panel border border-studio-border rounded px-3 py-2 pr-10 text-sm text-white placeholder:text-studio-muted/50 focus:outline-none focus:border-studio-accent"
          />
          <button
            type="button"
            onClick={() => setShowPassword((v) => !v)}
            aria-label={showPassword ? "Hide password" : "Show password"}
            className="absolute inset-y-0 right-0 px-3 text-studio-muted hover:text-white transition"
          >
            {showPassword ? <EyeOff className="w-4 h-4" /> : <Eye className="w-4 h-4" />}
          </button>
        </div>
      </div>

      <button
        type="submit"
        disabled={isSubmitting || !token || password.length < 8}
        className="w-full bg-studio-accent hover:bg-studio-accent-dark disabled:opacity-50 text-white font-medium text-sm py-2.5 rounded transition-colors"
      >
        {isSubmitting ? "Updating..." : "Reset password"}
      </button>

      <p className="text-xs text-studio-muted text-center">
        <Link href="/login" className="text-studio-accent hover:underline">
          Back to sign in
        </Link>
      </p>
    </form>
  );
}

export default function ResetPasswordPage() {
  return (
    <div className="min-h-screen bg-studio-bg flex items-center justify-center p-8">
      <div className="w-full max-w-sm space-y-8">
        <div className="flex items-center gap-3 justify-center">
          <div className="w-10 h-10 rounded-md bg-studio-panel border border-studio-border flex items-center justify-center text-studio-accent">
            <Film className="w-5 h-5" />
          </div>
          <span className="font-bold text-xl tracking-tight text-white">DramaStudio</span>
        </div>
        <Suspense>
          <ResetPasswordForm />
        </Suspense>
      </div>
    </div>
  );
}
