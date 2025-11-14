<script lang="ts">
	import '../app.css';
	import favicon from '$lib/assets/favicon.svg';
	import Header from '$components/Header.svelte';
	import type { Snippet } from 'svelte';

	interface LayoutProps {
		children: Snippet;
	}

	let { children }: LayoutProps = $props();
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
</svelte:head>

<div class="min-h-screen bg-base-200">
	<Header />

	{@render children()}
</div>
