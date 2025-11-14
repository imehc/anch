<script lang="ts">
	import Card from '$components/Card.svelte';
	import SectionHeader from '$components/SectionHeader.svelte';
	import type { PageData } from './$types';
	import * as m from '$lib/paraglide/messages';

	let { data }: { data: PageData } = $props();
</script>

<!-- Main Content -->
<main class="mx-auto max-w-7xl px-3 py-4 sm:px-4 sm:py-6 lg:px-6 lg:py-8">
	<!-- Hero Section -->
	<!-- <section class="mb-6 sm:mb-8 lg:mb-12">
		<h2 class="mb-1 text-2xl font-bold sm:mb-2 sm:text-3xl lg:text-4xl">工具与游戏</h2>
		<p class="text-sm text-base-content/70 sm:text-base lg:text-lg">浏览所有可用的工具和游戏</p>
	</section> -->

	<!-- Dynamic Sections -->
	{#each data.sections as section}
		<section class="mb-8 sm:mb-10 lg:mb-12">
			<SectionHeader
				title={section.title}
				iconBgColor={section.iconBgColor}
				iconColor={section.iconColor}
			>
				{#snippet icon()}
					{#await import(`$components/icons/${section.iconType}.svelte`) then { default: IconComponent }}
						<IconComponent class="size-5 sm:size-6" />
					{/await}
				{/snippet}
			</SectionHeader>

			<div class="grid grid-cols-2 gap-3 sm:gap-4 md:grid-cols-3 lg:grid-cols-4">
				{#each section.items as item}
					<Card
						href={item.href}
						title={item.title}
						description={item.description}
						iconBgColor={item.iconBgColor}
						iconColor={item.iconColor}
						category={item.category}
					>
						{#snippet icon()}
							{#if item.iconType === 'number' && 'iconContent' in item}
								<div class="text-xl font-bold {item.iconColor} sm:text-2xl">
									{item.iconContent}
								</div>
							{:else}
								{#await import(`$components/icons/${item.iconType}.svelte`) then { default: IconComponent }}
									<IconComponent class="size-6 sm:size-7 lg:size-8" />
								{/await}
							{/if}
						{/snippet}
					</Card>
				{:else}
					<div class="col-span-2 rounded-lg py-8 text-center md:col-span-3 lg:col-span-4">
						<p class="text-base-content/60">{m.category_empty()}</p>
					</div>
				{/each}
			</div>
		</section>
	{:else}
		<div class="flex min-h-[60vh] flex-col items-center justify-center text-center">
			<div class="mb-4 text-6xl opacity-20">🔍</div>
			<h3 class="mb-2 text-xl font-semibold sm:text-2xl">{m.no_results_title()}</h3>
			<p class="mb-6 text-base-content/70">
				{#if data.query}
					{m.no_results_message({ query: data.query })}
				{:else}
					{m.no_content_message()}
				{/if}
			</p>
			{#if data.query}
				<a href="/" class="btn btn-primary btn-sm sm:btn-md">
					<svg
						class="h-4 w-4 sm:h-5 sm:w-5"
						fill="none"
						stroke="currentColor"
						viewBox="0 0 24 24"
					>
						<path
							stroke-linecap="round"
							stroke-linejoin="round"
							stroke-width="2"
							d="M10 19l-7-7m0 0l7-7m-7 7h18"
						/>
					</svg>
					{m.back_home()}
				</a>
			{/if}
		</div>
	{/each}
</main>
