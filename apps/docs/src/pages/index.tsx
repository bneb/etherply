import type { ReactNode } from 'react';
import clsx from 'clsx';
import Link from '@docusaurus/Link';
import useDocusaurusContext from '@docusaurus/useDocusaurusContext';
import Layout from '@theme/Layout';
import HomepageFeatures from '@site/src/components/HomepageFeatures';
import Heading from '@theme/Heading';

import styles from './index.module.css';

function HomepageHeader() {
  const { siteConfig } = useDocusaurusContext();
  return (
    <header className={clsx('hero', styles.heroBanner)}>
      <div className="container">
        <h1 className="hero__title">
          {siteConfig.title}
        </h1>
        <p className="hero__subtitle" style={{ fontSize: '1.5rem', fontWeight: 300, opacity: 0.8 }}>
          The infrastructure of presence.
        </p>
        <div style={{ maxWidth: '600px', margin: '0 auto', padding: '2rem 0', lineHeight: '1.8' }}>
          <p>
            The old web was a library—static, archived, solitary. The new web is a conversation.
            EtherPly transforms dead interfaces into living rooms.
            Synchronize state across the globe in milliseconds, with the elegance of a local variable.
          </p>
        </div>
        <div className={styles.buttons}>
          <Link
            className="button button--secondary button--lg"
            to="/docs/intro"
            style={{ borderRadius: '0', textTransform: 'uppercase', padding: '1rem 2rem', letterSpacing: '1px' }}>
            Begin the Integration
          </Link>
        </div>
      </div>
    </header>
  );
}

export default function Home(): ReactNode {
  const { siteConfig } = useDocusaurusContext();
  return (
    <Layout
      title={`Hello from ${siteConfig.title}`}
      description="Description will go into a meta tag in <head />">
      <HomepageHeader />
      <main>
        <HomepageFeatures />
      </main>
    </Layout>
  );
}
