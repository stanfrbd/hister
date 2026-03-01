<script lang="ts">
  import { page } from '$app/stores';
  import { onMount } from 'svelte';
  import { Button } from '@hister/components/ui/button';
  import { Sun } from 'lucide-svelte';
  import "../style.css";

  let { children } = $props();

  function applyTheme(theme: string) {
    document.documentElement.setAttribute('data-theme', theme);
    if (theme === 'dark') {
      document.documentElement.classList.add('dark');
    } else {
      document.documentElement.classList.remove('dark');
    }
  }

  onMount(() => {
    const theme = localStorage.getItem('theme') ||
      (window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light');
    applyTheme(theme);
  });

  function toggleTheme() {
    const current = document.documentElement.getAttribute('data-theme');
    const next = current === 'dark' ? 'light' : 'dark';
    applyTheme(next);
    localStorage.setItem('theme', next);
  }
</script>

<header class="h-16 px-6 bg-page-bg border-b-[2px] border-border-brand flex items-center justify-between shadow-[4px_4px_0px_var(--hister-indigo)]">
  <h1 class="flex items-center gap-2">
    <img src="static/logo.png" alt="Hister logo" class="h-8 w-8" />
    <a data-sveltekit-reload href="./" class="font-outfit text-xl font-extrabold text-text-brand no-underline hover:underline">
      Hister
    </a>
  </h1>
  <nav class="flex items-center gap-6">
    <a
      class:underline={$page.url.pathname === new URL('history', $page.url).pathname}
      class="font-inter text-sm font-semibold text-text-brand-secondary hover:text-text-brand no-underline hover:underline"
      href="history"
    >
      History
    </a>
    <a
      class:underline={$page.url.pathname === new URL('rules', $page.url).pathname}
      class="font-inter text-sm font-semibold text-text-brand-secondary hover:text-text-brand no-underline hover:underline"
      href="rules"
    >
      Rules
    </a>
    <a
      class:underline={$page.url.pathname === new URL('add', $page.url).pathname}
      class="font-inter text-sm font-semibold text-text-brand-secondary hover:text-text-brand no-underline hover:underline"
      href="add"
    >
      Add
    </a>
  </nav>
  <Button
    variant="ghost"
    size="icon"
    class="text-text-brand-muted hover:text-hister-indigo transition-all hover:scale-110"
    title="Toggle theme"
    onclick={toggleTheme}
  >
    <Sun class="size-6" />
  </Button>
</header>

<main class="flex flex-col overflow-hidden">
  {@render children()}
</main>

<footer class="h-12 px-6 bg-page-bg border-t-[2px] border-border-brand flex items-center justify-center gap-4 text-sm">
  <a href="help" class="font-inter text-text-brand-secondary hover:text-hister-indigo no-underline hover:underline">Help</a>
  <span class="text-text-brand-muted">|</span>
  <a href="about" class="font-inter text-text-brand-secondary hover:text-hister-indigo no-underline hover:underline">About</a>
  <span class="text-text-brand-muted">|</span>
  <a href="api" class="font-inter text-text-brand-secondary hover:text-hister-indigo no-underline hover:underline">API</a>
  <span class="text-text-brand-muted">|</span>
  <a href="https://github.com/asciimoo/hister/" class="font-inter text-text-brand-secondary hover:text-hister-indigo no-underline hover:underline" target="_blank" rel="noopener">GitHub</a>
</footer>
