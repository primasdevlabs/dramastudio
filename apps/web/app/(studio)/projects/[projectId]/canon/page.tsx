"use client";

import { useState } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { ShieldCheck, Plus, CheckCircle2, AlertTriangle, XCircle } from "lucide-react";
import { Paper, Group, Stack, Title, Text, Badge, Button, Modal, TextInput, Table, ThemeIcon } from "@mantine/core";
import { notifications } from "@mantine/notifications";
import { api } from "@/lib/api/client";
import { StoryFact } from "@/lib/api/types";

export default function CanonPage({ params }: { params: { projectId: string } }) {
  const { projectId } = params;
  const queryClient = useQueryClient();

  const [subject, setSubject] = useState("");
  const [predicate, setPredicate] = useState("owns");
  const [object, setObject] = useState("");
  const [opened, setOpened] = useState(false);

  const { data: facts, isLoading } = useQuery({
    queryKey: ["canon-facts", projectId],
    queryFn: async () => {
      const res = await api.get<{ facts: StoryFact[] }>(`/v1/canon/facts?project_id=${projectId}`);
      return res.facts || [];
    },
  });

  const createMutation = useMutation({
    mutationFn: async () => {
      return api.post<StoryFact>("/v1/canon/facts", {
        subject,
        predicate,
        object,
        introduced: "Episode 1",
        valid_from: "Episode 1",
      });
    },
    onSuccess: (fact) => {
      queryClient.invalidateQueries({ queryKey: ["canon-facts", projectId] });
      setOpened(false);
      setSubject("");
      setObject("");
      notifications.show({
        title: "Fact Established",
        message: `Added canonical story fact: ${fact.subject || subject} ${fact.predicate || predicate} ${fact.object || object}`,
        color: "cyan",
      });
    },
  });

  return (
    <Stack gap="lg" className="max-w-7xl mx-auto">
      {/* Header */}
      <Paper p="xl" radius="xl" withBorder className="bg-studio-card border-studio-border shadow-xl">
        <Group justify="space-between" align="center">
          <Group gap="md">
            <ThemeIcon color="cyan" variant="light" size={44} radius="xl">
              <ShieldCheck size={22} />
            </ThemeIcon>
            <div>
              <Title order={3} c="white">
                Canon Context & Story Facts
              </Title>
              <Text size="xs" c="dimmed">
                Authoritative source of story truth and character knowledge state
              </Text>
            </div>
          </Group>

          <Button
            onClick={() => setOpened(true)}
            leftSection={<Plus size={16} />}
            variant="gradient"
            gradient={{ from: "cyan", to: "blue", deg: 90 }}
            color="cyan"
          >
            Add Canonical Fact
          </Button>
        </Group>
      </Paper>

      {/* Modal */}
      <Modal
        opened={opened}
        onClose={() => setOpened(false)}
        title={
          <Text fw={700} size="md" c="white">
            Establish Story Fact
          </Text>
        }
        centered
        overlayProps={{ backgroundOpacity: 0.7, blur: 8 }}
      >
        <Stack gap="md">
          <TextInput
            label="Subject (Character / Entity)"
            placeholder="e.g. Sarah"
            required
            value={subject}
            onChange={(e) => setSubject(e.currentTarget.value)}
            variant="filled"
          />

          <TextInput
            label="Predicate (Relation / Fact)"
            placeholder="e.g. discovered_secret"
            value={predicate}
            onChange={(e) => setPredicate(e.currentTarget.value)}
            variant="filled"
          />

          <TextInput
            label="Object (Target / Information)"
            placeholder="e.g. confidential_file"
            required
            value={object}
            onChange={(e) => setObject(e.currentTarget.value)}
            variant="filled"
          />

          <Group justify="flex-end" gap="sm" mt="md">
            <Button variant="subtle" color="gray" onClick={() => setOpened(false)}>
              Cancel
            </Button>
            <Button
              onClick={() => createMutation.mutate()}
              loading={createMutation.isPending}
              disabled={!subject || !object}
              color="cyan"
            >
              Establish Fact
            </Button>
          </Group>
        </Stack>
      </Modal>

      {/* Facts Table */}
      <Paper radius="xl" withBorder className="bg-studio-card border-studio-border overflow-hidden shadow-xl">
        <Table verticalSpacing="sm" horizontalSpacing="lg">
          <Table.Thead className="bg-studio-panel">
            <Table.Tr>
              <Table.Th c="dimmed">Subject</Table.Th>
              <Table.Th c="dimmed">Predicate</Table.Th>
              <Table.Th c="dimmed">Object</Table.Th>
              <Table.Th c="dimmed">Introduced In</Table.Th>
              <Table.Th c="dimmed">Status</Table.Th>
            </Table.Tr>
          </Table.Thead>
          <Table.Tbody>
            {isLoading ? (
              <Table.Tr>
                <Table.Td colSpan={5} align="center">
                  <Text size="xs" c="dimmed" fs="italic">
                    Loading canonical story facts...
                  </Text>
                </Table.Td>
              </Table.Tr>
            ) : facts && facts.length > 0 ? (
              facts.map((fact) => (
                <Table.Tr key={fact.id}>
                  <Table.Td fw={700} c="white">
                    {fact.subject}
                  </Table.Td>
                  <Table.Td className="font-semibold text-cyan-400">
                    {fact.predicate}
                  </Table.Td>
                  <Table.Td fw={600} c="white">
                    {fact.object}
                  </Table.Td>
                  <Table.Td c="dimmed">{fact.introduced || "Episode 1"}</Table.Td>
                  <Table.Td>
                    {fact.status === "canonical" ? (
                      <Badge color="emerald" variant="light" size="sm" leftSection={<CheckCircle2 size={12} />}>
                        Canonical
                      </Badge>
                    ) : fact.status === "disputed" ? (
                      <Badge color="amber" variant="light" size="sm" leftSection={<AlertTriangle size={12} />}>
                        Disputed
                      </Badge>
                    ) : (
                      <Badge color="red" variant="light" size="sm" leftSection={<XCircle size={12} />}>
                        Retconned
                      </Badge>
                    )}
                  </Table.Td>
                </Table.Tr>
              ))
            ) : (
              <Table.Tr>
                <Table.Td colSpan={5} align="center">
                  <Text size="xs" c="dimmed" fs="italic">
                    No canonical facts established yet.
                  </Text>
                </Table.Td>
              </Table.Tr>
            )}
          </Table.Tbody>
        </Table>
      </Paper>
    </Stack>
  );
}
