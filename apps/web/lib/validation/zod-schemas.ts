import { z } from 'zod'

export const StoryboardPanelSchema = z.object({
  id: z.string().uuid(),
  shotNumber: z.number().int().positive(),
  prompt: z.string().min(5),
  durationSec: z.number().positive(),
  characterReferences: z.array(z.string()),
  cameraAngle: z.string(),
  capability: z.literal('storyboard'),
})

export type StoryboardPanelInput = z.infer<typeof StoryboardPanelSchema>

export const GenerationJobSchema = z.object({
  id: z.string(),
  capability: z.string(),
  modelId: z.string(),
  status: z.enum(['pending', 'executing', 'completed', 'failed']),
  resultUrl: z.string().url().optional(),
})

export type GenerationJobInput = z.infer<typeof GenerationJobSchema>
