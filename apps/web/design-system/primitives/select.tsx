"use client";

import { Select as MantineSelect, SelectProps } from "@mantine/core";
import { forwardRef } from "react";

export const Select = forwardRef<HTMLInputElement, SelectProps>(
  (props, ref) => {
    return <MantineSelect ref={ref} variant="filled" radius="md" {...props} />;
  }
);
Select.displayName = "StudioSelect";
