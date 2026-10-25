import AgentConsole from "../AgentConsole";
import { HERO, LINKS, SITE } from "../../lib/profile";

export default function Hero() {
  return (
    <section className="wrap hero" id="top">
      <div>
        <div className="eyebrow">{HERO.eyebrow}</div>
        <h1>
          {SITE.name}
          <br />
          <span className="grad">
            {HERO.headline}
            <br />
            {HERO.headlineAccent}
          </span>
        </h1>
        <div className="role">{HERO.role}</div>
        <p className="lede">{HERO.lede}</p>
        <div className="cta">
          <a className="btn primary" href="#agent">
            Ask my AI agent
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
              <path d="M5 12h14M13 6l6 6-6 6" />
            </svg>
          </a>
          {SITE.resumeUrl && (
            <a className="btn" href={SITE.resumeUrl} download>
              Résumé
            </a>
          )}
          {LINKS.github && (
            <a className="btn" href={LINKS.github} target="_blank" rel="noopener noreferrer">
              GitHub
            </a>
          )}
          {LINKS.linkedin && (
            <a className="btn" href={LINKS.linkedin} target="_blank" rel="noopener noreferrer">
              LinkedIn
            </a>
          )}
        </div>
      </div>

      <AgentConsole />
    </section>
  );
}
