"use client";

import { Modal as MantineModal, ModalProps } from "@mantine/core";

export function Dialog(props: ModalProps) {
  return (
    <MantineModal
      centered
      radius="xl"
      overlayProps={{ backgroundOpacity: 0.7, blur: 8 }}
      {...props}
    />
  );
}
