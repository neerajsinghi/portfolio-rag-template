// ─────────────────────────────────────────────────────────────────────────────
//  MAKE IT YOURS: this is the ONLY file you need to edit for the website.
//
//  Everything below is a PLACEHOLDER persona ("Alex Doe"). Replace each value
//  with your own. Page copy, SEO metadata, JSON-LD, the Open Graph image,
//  robots.txt, sitemap.xml and llms.txt are all generated from this file.
//
//  The AI agent's knowledge lives separately in backend/corpus/ (profile.json
//  plus Markdown documents) — see backend/corpus/README.md.
// ─────────────────────────────────────────────────────────────────────────────

export const SITE = {
  /** Full name, shown in the hero, title, footer and structured data. */
  name: "Alex Doe",
  givenName: "Alex",
  familyName: "Doe",
  /** Short handle used in the nav brand and console title. */
  handle: "alexdoe",
  /** Public URL of the deployed site (no trailing slash). NEXT_PUBLIC_SITE_URL overrides it. */
  url: "https://example.com",
  /** One-line role, used as the job title in metadata. */
  jobTitle: "Senior Software Engineer",
  location: "Springfield, Example Country",
  /** Short status shown in the nav pill. Set to "" to hide it. */
  status: "Open to Senior / Staff roles",
  /** <title> and social-card title. */
  title: "Alex Doe | Senior Software Engineer",
  /** Meta description (aim for 120–160 characters). */
  description:
    "Alex Doe is a placeholder senior software engineer. Ask the AI agent on this site about experience, skills and projects.",
  /** Keywords meta tag (Bing and some people-search tools read it). */
  keywords: ["Alex Doe", "software engineer", "backend engineer", "portfolio"],
  /** Open Graph locale, e.g. en_US, en_GB, en_IN. */
  locale: "en_US",
  /** Optional résumé file placed in frontend/public/. Set to "" to hide the buttons. */
  resumeUrl: "",
};

export const HERO = {
  eyebrow: "Backend · Cloud · AI",
  headline: "Senior software engineer",
  headlineAccent: "building reliable systems.",
  role: "Placeholder persona · ~8 years shipping backend services, developer tools and cloud infrastructure.",
  lede:
    "Replace this paragraph with a two-sentence pitch. Visitors can ask the AI agent anything about your work — it answers from your documents using retrieval and tools.",
};

/** Paragraphs for the About section. */
export const ABOUT: string[] = [
  "This is placeholder copy for the About section. Write two or three short paragraphs in the first person about what you build, what you care about, and what kind of work you want next.",
  "Keep it concrete: mention a system you built, a problem you solved, or a team you helped. The agent's knowledge base in backend/corpus/ can go deeper than this page.",
];

export const FACTS: [string, string][] = [
  ["experience", "~8 years"],
  ["core", "Go · TypeScript · Python"],
  ["cloud", "Docker · Kubernetes · Terraform"],
  ["focus now", "AI · RAG · agents"],
  ["based in", SITE.location],
  ["open to", "Senior / Staff · remote"],
];

export const STACK: { title: string; tags: string[] }[] = [
  { title: "Languages", tags: ["Go", "TypeScript", "Python"] },
  { title: "Backend", tags: ["REST", "gRPC", "Microservices", "PostgreSQL", "Redis"] },
  { title: "Cloud & Infra", tags: ["Docker", "Kubernetes", "Terraform", "CI/CD"] },
  { title: "AI / LLM", tags: ["RAG", "Tool-use agents", "Prompt design"] },
];

export type Job = { role: string; co: string; when: string; dim?: boolean; pts: string[] };

export const EXPERIENCE: Job[] = [
  {
    role: "Senior Software Engineer",
    co: "Example Corp",
    when: "2022 — Present",
    pts: [
      "Placeholder: leads a small platform team and owns the service template used across the company.",
      "Placeholder: cut deployment time from an hour to ten minutes with a new CI/CD pipeline.",
    ],
  },
  {
    role: "Software Engineer",
    co: "Sample Startup Inc.",
    when: "2018 — 2022",
    pts: ["Placeholder: built the public REST API and an event-driven data pipeline from prototype to production."],
  },
  {
    role: "Junior Developer",
    co: "First Job LLC",
    when: "2016 — 2018",
    dim: true,
    pts: ["Placeholder: earlier roles can be dimmed with `dim: true`."],
  },
];

export const PROJECTS: { tag: string; name: string; desc: string; stack: string[]; url?: string }[] = [
  {
    tag: "Open source",
    name: "Example Project One",
    desc: "Placeholder: a short description of a project you are proud of and the impact it had.",
    stack: ["Go", "PostgreSQL", "Docker"],
  },
  {
    tag: "AI · RAG",
    name: "Portfolio Agent (this site)",
    desc: "A portfolio with an AI agent that answers questions about me using retrieval over my documents and a tool-use loop.",
    stack: ["Go", "Next.js", "RAG", "SSE"],
  },
];

/** Rendered on the page and as FAQPage JSON-LD. Answer the questions people search for. */
export const FAQ: { q: string; a: string }[] = [
  {
    q: "Who is Alex Doe?",
    a: "Alex Doe is the placeholder persona that ships with this template. Replace this answer with a one-paragraph summary of who you are and what you do.",
  },
  {
    q: "What roles is Alex Doe open to?",
    a: "Placeholder: Senior and Staff engineering roles, remote-friendly. Replace with your own availability.",
  },
  {
    q: "How can I contact Alex Doe?",
    a: "Placeholder: email alex@example.com or use the links in the contact section.",
  },
];

/** Suggested questions shown under the agent console. */
export const CHIPS = [
  "Give me a quick overview of Alex",
  "What is Alex's tech stack?",
  "Walk me through Alex's experience",
  "Is Alex a fit for a Staff backend role?",
  "How can I contact Alex?",
];

/** First message shown in the agent console (Markdown). */
export const GREETING = `Hi — I'm **${SITE.givenName}'s portfolio agent**. I answer from ${SITE.givenName}'s documents and GitHub using retrieval + tools. Ask about experience, skills, projects, or fit for a role.`;

/** Contact links. Set a value to "" to hide that card. */
export const LINKS = {
  email: "alex@example.com",
  linkedin: "https://www.linkedin.com/in/your-handle",
  github: "https://github.com/your-username",
};

/** Other profiles for JSON-LD `sameAs` (e.g. Mastodon, personal blog). */
export const SAME_AS: string[] = [LINKS.github, LINKS.linkedin].filter(Boolean);

// ── Runtime config (from env; see frontend/.env.example) ────────────────────
export const API_BASE = process.env.NEXT_PUBLIC_API_BASE || "http://localhost:8080";
export const siteUrl = (process.env.NEXT_PUBLIC_SITE_URL || SITE.url).replace(/\/$/, "");
/** Google Analytics 4 measurement ID. Analytics are disabled when empty. */
export const GA_ID = process.env.NEXT_PUBLIC_GA_ID || "";
