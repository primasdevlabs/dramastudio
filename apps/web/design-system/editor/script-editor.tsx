"use client";

import { useEditor, EditorContent } from "@tiptap/react";
import StarterKit from "@tiptap/starter-kit";
import { Paper, Group, ActionIcon, Tooltip } from "@mantine/core";
import { Bold, Italic, List, Heading1, Heading2 } from "lucide-react";

export function ScriptEditor({
  content = "<p>INT. COFFEE SHOP - DAY</p><p>SARAH sits across from MARK. The tension is palpable.</p>",
  onChange,
}: {
  content?: string;
  onChange?: (html: string) => void;
}) {
  const editor = useEditor({
    extensions: [StarterKit],
    content,
    onUpdate: ({ editor }: { editor: any }) => {
      if (onChange) onChange(editor.getHTML());
    },
  });

  if (!editor) return null;

  return (
    <Paper p="sm" radius="lg" withBorder className="bg-studio-panel border-studio-border">
      <Group gap="xs" mb="xs" pb="xs" className="border-b border-studio-border">
        <Tooltip label="Bold" withArrow>
          <ActionIcon
            variant={editor.isActive("bold") ? "filled" : "subtle"}
            color="cyan"
            size="sm"
            onClick={() => editor.chain().focus().toggleBold().run()}
          >
            <Bold size={14} />
          </ActionIcon>
        </Tooltip>

        <Tooltip label="Italic" withArrow>
          <ActionIcon
            variant={editor.isActive("italic") ? "filled" : "subtle"}
            color="cyan"
            size="sm"
            onClick={() => editor.chain().focus().toggleItalic().run()}
          >
            <Italic size={14} />
          </ActionIcon>
        </Tooltip>

        <Tooltip label="Scene Heading (H1)" withArrow>
          <ActionIcon
            variant={editor.isActive("heading", { level: 1 }) ? "filled" : "subtle"}
            color="cyan"
            size="sm"
            onClick={() => editor.chain().focus().toggleHeading({ level: 1 }).run()}
          >
            <Heading1 size={14} />
          </ActionIcon>
        </Tooltip>

        <Tooltip label="Action Heading (H2)" withArrow>
          <ActionIcon
            variant={editor.isActive("heading", { level: 2 }) ? "filled" : "subtle"}
            color="cyan"
            size="sm"
            onClick={() => editor.chain().focus().toggleHeading({ level: 2 }).run()}
          >
            <Heading2 size={14} />
          </ActionIcon>
        </Tooltip>

        <Tooltip label="Bullet List" withArrow>
          <ActionIcon
            variant={editor.isActive("bulletList") ? "filled" : "subtle"}
            color="cyan"
            size="sm"
            onClick={() => editor.chain().focus().toggleBulletList().run()}
          >
            <List size={14} />
          </ActionIcon>
        </Tooltip>
      </Group>

      <div className="text-sm text-studio-text min-h-[160px] p-2 focus:outline-none">
        <EditorContent editor={editor} />
      </div>
    </Paper>
  );
}
