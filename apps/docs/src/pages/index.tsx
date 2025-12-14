
import React from 'react';
import clsx from 'clsx';
import Layout from '@theme/Layout';
import Link from '@docusaurus/Link';
import useDocusaurusContext from '@docusaurus/useDocusaurusContext';
import styles from './index.module.css';

function HomepageHeader() {
  const { siteConfig } = useDocusaurusContext();
  return (
    <header className={clsx('hero', styles.heroBanner)} style={{ padding: '4rem 0', background: 'transparent' }}>
      <div className="container">
        <div style={{ display: 'flex', flexDirection: 'column', alignItems: 'center', textAlign: 'center', marginBottom: '3rem' }}>
          <h1 className="hero__title" style={{ fontSize: '3.5rem', fontWeight: 800, letterSpacing: '-0.02em', marginBottom: '1.5rem' }}>
            <span style={{ background: 'linear-gradient(to right, #60a5fa, #3b82f6)', WebkitBackgroundClip: 'text', WebkitTextFillColor: 'transparent' }}>
              Global Sync
            </span> in &lt;50ms.
          </h1>
          <p className="hero__subtitle" style={{ fontSize: '1.25rem', color: '#a1a1aa', maxWidth: '600px', margin: '0 auto 2rem' }}>
            EtherPly is the open-source sync engine for modern applications.
            Stop building WebSockets from scratch.
          </p>

          <div style={{ display: 'flex', gap: '1rem', justifyContent: 'center' }}>
            <Link
              className="button button--primary button--lg"
              to="http://localhost:3000/login"
              style={{ height: '3rem', display: 'flex', alignItems: 'center' }}>
              Get Started →
            </Link>
            <Link
              className="button button--secondary button--lg"
              to="/docs/intro"
              style={{ height: '3rem', display: 'flex', alignItems: 'center', background: 'rgba(255,255,255,0.1)', border: 'none', color: '#fff' }}>
              Read Documentation
            </Link>
          </div>
        </div>

        {/* Feature Grid Visualization */}
        <div style={{
          display: 'grid',
          gridTemplateColumns: 'repeat(auto-fit, minmax(300px, 1fr))',
          gap: '2rem',
          maxWidth: '1200px',
          margin: '0 auto',
          padding: '0 1rem'
        }}>
          {/* Card 1: Collaborative Editor */}
          <div className="feature-card" style={{ borderRadius: '12px', overflow: 'hidden', border: '1px solid rgba(255,255,255,0.1)', background: 'rgba(255,255,255,0.02)' }}>
            <div style={{ padding: '0.75rem 1rem', borderBottom: '1px solid rgba(255,255,255,0.05)', fontSize: '0.9rem', color: '#a1a1aa', fontWeight: 600 }}>
              Multiplayer Primitives
            </div>
            <img src="/docs/img/view-editor.png" alt="Collaborative Editor" style={{ display: 'block', width: '100%', height: 'auto' }} />
          </div>

          {/* Card 2: Analytics */}
          <div className="feature-card" style={{ borderRadius: '12px', overflow: 'hidden', border: '1px solid rgba(255,255,255,0.1)', background: 'rgba(255,255,255,0.02)' }}>
            <div style={{ padding: '0.75rem 1rem', borderBottom: '1px solid rgba(255,255,255,0.05)', fontSize: '0.9rem', color: '#a1a1aa', fontWeight: 600 }}>
              Real-time Insights
            </div>
            <img src="/docs/img/view-analytics.png" alt="Real-time Analytics" style={{ display: 'block', width: '100%', height: 'auto' }} />
          </div>

          {/* Card 3: Topology */}
          <div className="feature-card" style={{ borderRadius: '12px', overflow: 'hidden', border: '1px solid rgba(255,255,255,0.1)', background: 'rgba(255,255,255,0.02)' }}>
            <div style={{ padding: '0.75rem 1rem', borderBottom: '1px solid rgba(255,255,255,0.05)', fontSize: '0.9rem', color: '#a1a1aa', fontWeight: 600 }}>
              Global Mesh
            </div>
            <img src="/docs/img/view-topology.png" alt="Global Infrastructure" style={{ display: 'block', width: '100%', height: 'auto' }} />
          </div>
        </div>

        {/* Quick Install Code Block */}
        <div style={{ maxWidth: '600px', margin: '4rem auto 0' }}>
          <div style={{
            background: '#18181b',
            borderRadius: '8px',
            border: '1px solid #27272a',
            overflow: 'hidden'
          }}>
            <div style={{
              padding: '0.5rem 1rem',
              background: '#27272a',
              color: '#a1a1aa',
              fontSize: '0.8rem',
              fontWeight: 600
            }}>
              INSTALLATION
            </div>
            <div style={{ padding: '1.5rem', fontFamily: 'monospace', color: '#e4e4e7', display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
              <span>npm install @etherply/sdk</span>
              <button
                onClick={() => navigator.clipboard.writeText('npm install @etherply/sdk')}
                style={{ background: 'none', border: 'none', cursor: 'pointer', color: '#71717a' }}
                title="Copy"
              >
                📋
              </button>
            </div>
          </div>
        </div>

      </div>
    </header>
  );
}

export default function Home(): React.JSX.Element {
  const { siteConfig } = useDocusaurusContext();
  return (
    <Layout
      title={`EtherPly Documentation`}
      description="Real-time sync engine for professional applications">
      <HomepageHeader />
    </Layout>
  );
}
