"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import { Film, Eye, EyeOff } from "lucide-react";
import { api, setAuthToken } from "@/lib/api/client";
import { GoogleSignIn } from "@/components/google-sign-in";

export default function RegisterPage() {
  const router = useRouter();
  const [orgName, setOrgName] = useState("");
  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [showPassword, setShowPassword] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [isSubmitting, setIsSubmitting] = useState(false);

  const handleRegister = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!orgName || !name || !email || password.length < 8) return;
    setIsSubmitting(true);
    setError(null);
    try {
      await api.post("/v1/auth/register", {
        org_name: orgName,
        email,
        name,
        password,
      });
      // Register returns the user, not a token — log in immediately after.
      const res = await api.post<{ access_token: string }>("/v1/auth/login", {
        email,
        password,
      });
      setAuthToken(res.access_token);
      router.push("/projects");
    } catch (err) {
      setError(err instanceof Error ? err.message : "Registration failed.");
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

        <form
          onSubmit={handleRegister}
          className="bg-studio-card border border-studio-border rounded-md p-6 space-y-4"
        >
          <div>
            <h1 className="text-lg font-bold text-white">Create your studio</h1>
            <p className="text-xs text-studio-muted mt-1">
              Registering creates a new organization with you as its owner.
            </p>
          </div>

          {error && (
            <div className="bg-red-500/10 border border-red-500/30 text-red-400 rounded px-3 py-2 text-xs">
              {error}
            </div>
          )}

          <div className="space-y-1">
            <label className="block text-xs font-semibold text-studio-muted">Organization name</label>
            <input
              type="text"
              value={orgName}
              onChange={(e) => setOrgName(e.target.value)}
              placeholder="Acme Studios"
              autoComplete="organization"
              className="w-full bg-studio-panel border border-studio-border rounded px-3 py-2 text-sm text-white focus:outline-none focus:border-studio-accent"
            />
          </div>

          <div className="space-y-1">
            <label className="block text-xs font-semibold text-studio-muted">Your name</label>
            <input
              type="text"
              value={name}
              onChange={(e) => setName(e.target.value)}
              placeholder="Ridley Scott"
              autoComplete="name"
              className="w-full bg-studio-panel border border-studio-border rounded px-3 py-2 text-sm text-white focus:outline-none focus:border-studio-accent"
            />
          </div>

          <div className="space-y-1">
            <label className="block text-xs font-semibold text-studio-muted">Email</label>
            <input
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              placeholder="you@studio.com"
              autoComplete="email"
              className="w-full bg-studio-panel border border-studio-border rounded px-3 py-2 text-sm text-white focus:outline-none focus:border-studio-accent"
            />
          </div>

          <div className="space-y-1">
            <label className="block text-xs font-semibold text-studio-muted">Password</label>
            <div className="relative">
              <input
                type={showPassword ? "text" : "password"}
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                placeholder="Minimum 8 characters"
                autoComplete="new-password"
                className="w-full bg-studio-panel border border-studio-border rounded px-3 py-2 pr-10 text-sm text-white focus:outline-none focus:border-studio-accent"
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
            <p className="text-[10px] text-studio-muted">Minimum 8 characters.</p>
          </div>

          <button
            type="submit"
            disabled={isSubmitting || !orgName || !name || !email || password.length < 8}
            className="w-full bg-studio-accent hover:bg-studio-accent-dark disabled:opacity-50 text-white font-medium text-sm py-2.5 rounded transition-colors"
          >
            {isSubmitting ? "Creating account..." : "Create account"}
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

          <GoogleSignIn orgName={orgName || undefined} />

          <p className="text-xs text-studio-muted text-center">
            Already have an account?{" "}
            <Link href="/login" className="text-studio-accent hover:underline">
              Sign in
            </Link>
          </p>
        </form>
      </div>
    </div>
  );
}
