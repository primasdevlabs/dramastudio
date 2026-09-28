import { z } from "zod";

export const createProjectSchema = z.object({
  name: z.string().min(1, "Production title is required").max(200),
  description: z.string().max(2000).optional(),
  genre: z.string().min(1),
  language: z.string().min(2).max(5),
  mode: z.enum(["monitored", "autonomous"]),
});

export type CreateProjectInput = z.infer<typeof createProjectSchema>;
