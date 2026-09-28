"use client";

import { createContext, useContext } from "react";
import { StudioShell } from "@/components/layout/studio-shell";

const ProjectIdContext = createContext<string | undefined>(undefined);

export function useProjectId() {
  return useContext(ProjectIdContext);
}

export function StudioShellWrapper({
  children,
  projectId,
}: {
  children: React.ReactNode;
  projectId?: string;
}) {
  return (
    <ProjectIdContext.Provider value={projectId}>
      <StudioShell>{children}</StudioShell>
    </ProjectIdContext.Provider>
  );
}
