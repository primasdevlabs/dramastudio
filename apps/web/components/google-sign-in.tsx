"use client";

import { useEffect, useRef, useState } from "react";
import { useRouter } from "next/navigation";
import { api, setAuthToken } from "@/lib/api/client";

const GOOGLE_CLIENT_ID = process.env.NEXT_PUBLIC_GOOGLE_CLIENT_ID || "";
const GSI_SRC = "https://accounts.google.com/gsi/client";

declare global {
  interface Window {
    google?: {
      accounts: {
        id: {
          initialize: (opts: {
            client_id: string;
            callback: (resp: { credential: string }) => void;
          }) => void;
          renderButton: (
            el: HTMLElement,
            opts: Record<string, string | number>
          ) => void;
        };
      };
    };
  }
}

// Renders the official Google sign-in button. The ID token is verified
// server-side; first-time users get an org + account provisioned
// automatically (auto account detection). Hidden when unconfigured.
export function GoogleSignIn({ orgName }: { orgName?: string }) {
  const router = useRouter();
  const buttonRef = useRef<HTMLDivElement>(null);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!GOOGLE_CLIENT_ID) return;

    const onCredential = async (resp: { credential: string }) => {
      try {
        const res = await api.post<{ access_token: string }>("/v1/auth/sso", {
          provider: "google",
          id_token: resp.credential,
          org_name: orgName,
        });
        setAuthToken(res.access_token);
        router.push("/projects");
      } catch {
        setError("Google sign-in failed. Try email instead.");
      }
    };

    const init = () => {
      window.google?.accounts.id.initialize({
        client_id: GOOGLE_CLIENT_ID,
        callback: onCredential,
      });
      if (buttonRef.current) {
        window.google?.accounts.id.renderButton(buttonRef.current, {
          theme: "filled_black",
          size: "large",
          width: "320",
          shape: "pill",
          text: "continue_with",
        });
      }
    };

    if (window.google) {
      init();
      return;
    }
    const script = document.createElement("script");
    script.src = GSI_SRC;
    script.async = true;
    script.defer = true;
    script.onload = init;
    document.head.appendChild(script);
  }, [orgName, router]);

  if (!GOOGLE_CLIENT_ID) return null;

  return (
    <div className="space-y-2">
      <div ref={buttonRef} className="flex justify-center" />
      {error && <p className="text-xs text-red-400 text-center">{error}</p>}
    </div>
  );
}
