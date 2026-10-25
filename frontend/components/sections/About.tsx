import { ABOUT, FACTS } from "../../lib/profile";

export default function About() {
  return (
    <section className="wrap" id="about">
      <div className="sec-head">
        <span className="sec-no">01</span>
        <h2>About</h2>
      </div>
      <div className="about-grid">
        <div className="about">
          {ABOUT.map((p, i) => (
            <p key={i}>{p}</p>
          ))}
        </div>
        <div className="facts">
          {FACTS.map(([k, v]) => (
            <div className="fact" key={k}>
              <span className="k">{k}</span>
              <span className="v">{v}</span>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}
