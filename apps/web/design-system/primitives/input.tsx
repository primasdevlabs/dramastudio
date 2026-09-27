"use client";

import { TextInput as MantineTextInput, Textarea as MantineTextarea, TextInputProps, TextareaProps } from "@mantine/core";
import { forwardRef } from "react";

export const Input = forwardRef<HTMLInputElement, TextInputProps>(
  (props, ref) => {
    return <MantineTextInput ref={ref} variant="filled" radius="md" {...props} />;
  }
);
Input.displayName = "StudioInput";

export const Textarea = forwardRef<HTMLTextAreaElement, TextareaProps>(
  (props, ref) => {
    return <MantineTextarea ref={ref} variant="filled" radius="md" {...props} />;
  }
);
Textarea.displayName = "StudioTextarea";
