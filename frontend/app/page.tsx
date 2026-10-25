import type { Metadata } from "next";
import Topology from "../components/Topology";
import Hero from "../components/sections/Hero";
import About from "../components/sections/About";
import Stack from "../components/sections/Stack";
import Experience from "../components/sections/Experience";
import Projects from "../components/sections/Projects";
import FAQ from "../components/sections/FAQ";
import Contact from "../components/sections/Contact";
import { FAQ as FAQ_ITEMS, SAME_AS, SITE, STACK, LINKS, siteUrl } from "../lib/profile";

export const metadata: Metadata = {
  keywords: SITE.keywords,
};

const profile = {
  "@context": "https://schema.org",
  "@type": "ProfilePage",
  "@id": `${siteUrl}/#profile`,
  url: siteUrl,
  name: SITE.title,
  isPartOf: { "@id": `${siteUrl}/#website` },
  mainEntity: {
    "@type": "Person",
    "@id": `${siteUrl}/#person`,
    name: SITE.name,
    givenName: SITE.givenName,
    familyName: SITE.familyName,
    url: siteUrl,
    image: `${siteUrl}/opengraph-image`,
    jobTitle: SITE.jobTitle,
    description: SITE.description,
    ...(LINKS.email ? { email: `mailto:${LINKS.email}` } : {}),
    homeLocation: { "@type": "Place", name: SITE.location },
    sameAs: SAME_AS,
    knowsAbout: STACK.flatMap((group) => group.tags),
  },
};

const faqData = {
  "@context": "https://schema.org",
  "@type": "FAQPage",
  "@id": `${siteUrl}/#faq`,
  about: { "@id": `${siteUrl}/#person` },
  mainEntity: FAQ_ITEMS.map((item) => ({
    "@type": "Question",
    name: item.q,
    acceptedAnswer: { "@type": "Answer", text: item.a },
  })),
};

export default function Page() {
  return (
    <>
      <Topology />
      <div className="veil" />

      <header>
        <div className="wrap nav">
          <div className="brand">
            <span className="blink">▸</span> <b>{SITE.handle}</b>
            <span className="blink">_</span>
          </div>
          <nav className="nav-links">
            <a href="#about">About</a>
            <a href="#stack">Stack</a>
            <a href="#work">Work</a>
            <a href="#projects">Projects</a>
            <a href="#agent">Ask AI</a>
            {SITE.status && (
              <span className="pill">
                <span className="dot" />
                {SITE.status}
              </span>
            )}
          </nav>
        </div>
      </header>

      <main>
        <Hero />
        <About />
        <Stack />
        <Experience />
        <Projects />
        <FAQ />
        <Contact />
      </main>

      <footer>
        <div className="wrap">
          <span>
            © {new Date().getFullYear()} {SITE.name} · {SITE.location}
          </span>
          <span>RAG + tool-use agent · Go · Next.js</span>
        </div>
      </footer>
      <script
        type="application/ld+json"
        dangerouslySetInnerHTML={{ __html: JSON.stringify(profile).replace(/</g, "\\u003c") }}
      />
      <script
        type="application/ld+json"
        dangerouslySetInnerHTML={{ __html: JSON.stringify(faqData).replace(/</g, "\\u003c") }}
      />
    </>
  );
}
