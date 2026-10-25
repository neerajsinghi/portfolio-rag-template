import { EXPERIENCE, FAQ, LINKS, SITE, STACK, siteUrl } from "../../lib/profile";

// llms.txt (https://llmstxt.org): a plain-markdown summary that AI assistants
// and their crawlers can read without parsing the interactive homepage.
// Built from the same profile data as the site so the two never drift.
export const dynamic = "force-static";

export async function GET() {
  const lines = [
    `# ${SITE.name}`,
    "",
    `> ${SITE.description}`,
    "",
    "## Profile",
    "",
    `- [Homepage](${siteUrl}/): experience, stack, projects and contact`,
    ...(SITE.resumeUrl ? [`- [Résumé](${new URL(SITE.resumeUrl, siteUrl).toString()})`] : []),
    ...(LINKS.linkedin ? [`- [LinkedIn](${LINKS.linkedin})`] : []),
    ...(LINKS.github ? [`- [GitHub](${LINKS.github})`] : []),
    ...(LINKS.email ? [`- Email: ${LINKS.email}`] : []),
    "",
    "## Experience",
    "",
    ...EXPERIENCE.map((job) => `- ${job.role}, ${job.co} (${job.when.replace(/\s+—\s+/, " to ")}): ${job.pts.join(" ")}`),
    "",
    "## Skills",
    "",
    ...STACK.map((group) => `- ${group.title}: ${group.tags.join(", ")}`),
    "",
    "## FAQ",
    "",
    ...FAQ.flatMap((item) => [`### ${item.q}`, "", item.a, ""]),
  ];

  return new Response(lines.join("\n"), {
    headers: { "Content-Type": "text/markdown; charset=utf-8" },
  });
}
