"use client";

import { Suspense, useState } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import Link from "next/link";
import { Film, Eye, EyeOff } from "lucide-react";
import { api, setAuthToken } from "@/lib/api/client";
import { GoogleSignIn } from "@/components/google-sign-in";

function LoginForm() {
  const router = useRouter();
  const params = useSearchParams();
  const justReset = params.get("reset") === "1";
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [showPassword, setShowPassword] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [isSubmitting, setIsSubmitting] = useState(false);

  const handleLogin = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!email || !password) return;
    setIsSubmitting(true);
    setError(null);
    try {
      const res = await api.post<{ access_token: string }>("/v1/auth/login", {
        email,
        password,
      });
      setAuthToken(res.access_token);
      router.push("/projects");
    } catch (err) {
      const status = (err as { status?: number }).status;
      setError(
        status === 401
          ? "Invalid email or password."
          : err instanceof Error
            ? err.message
            : "Sign-in failed. Try again."
      );
      setIsSubmitting(false);
    }
  };

  return (
    <form
      onSubmit={handleLogin}
      className="bg-studio-card border border-studio-border rounded-md p-6 space-y-4"
    >
      <div>
        <h1 className="text-lg font-bold text-white">Sign in to the studio</h1>
        <p className="text-xs text-studio-muted mt-1">
          Use your organization account credentials.
        </p>
      </div>

      {justReset && (
        <div className="bg-emerald-500/10 border border-emerald-500/30 text-emerald-400 rounded px-3 py-2 text-xs">
          Password updated — sign in with your new password.
        </div>
      )}

      {error && (
        <div className="bg-red-500/10 border border-red-500/30 text-red-400 rounded px-3 py-2 text-xs">
          {error}
        </div>
      )}

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

      <div className="space-y-1">
        <div className="flex items-center justify-between">
          <label className="block text-xs font-semibold text-studio-muted">Password</label>
          <Link
            href="/forgot-password"
            className="text-[10px] text-studio-accent hover:underline"
          >
            Forgot password?
          </Link>
        </div>
        <div className="relative">
          <input
            type={showPassword ? "text" : "password"}
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            placeholder="Your account password"
            autoComplete="current-password"
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
        disabled={isSubmitting || !email || !password}
        className="w-full bg-studio-accent hover:bg-studio-accent-dark disabled:opacity-50 text-white font-medium text-sm py-2.5 rounded transition-colors"
      >
        {isSubmitting ? "Signing in..." : "Sign in"}
      </button>

      <div className="relative">
        <div className="absolute inset-0 flex items-center">
          <div className="w-full border-t border-studio-border" />
        </div>
        <div className="relative flex justify-center">
          <span className="bg-studio-card px-2 text-[10px] text-studio-muted uppercase tracking-wider">
            or
          </span>
        </div>
      </div>

      <GoogleSignIn />

      <p className="text-xs text-studio-muted text-center">
        No account yet?{" "}
        <Link href="/register" className="text-studio-accent hover:underline">
          Create one
        </Link>
        {" · "}
        <Link href="/" className="hover:text-white transition">
          Back to home
        </Link>
      </p>
    </form>
  );
}

export default function LoginPage() {
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
          <LoginForm />
        </Suspense>
      </div>
    </div>
  );
}
