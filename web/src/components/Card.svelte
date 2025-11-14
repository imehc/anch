<script lang="ts">
	import type { Snippet } from 'svelte';

	interface CardProps {
		href: string;
		title: string;
		description: string;
		iconBgColor: string;
		iconColor?: string;
		icon: Snippet;
		category: string;
	}

	let {
		href,
		title,
		description,
		iconBgColor,
		iconColor = '',
		icon,
		category
	}: CardProps = $props();
</script>

<a {href} class="card bg-base-100 shadow-md transition-all duration-300 hover:shadow-xl hover:-translate-y-1">
	<div class="card-body aspect-video p-3 sm:p-4 lg:p-5">
		<div class="relative flex flex-col h-full gap-3">
			<!-- Category Badge - Top Right -->
			<div class="absolute -right-1 -top-1 sm:right-0 sm:top-0">
				<span class="badge badge-sm bg-base-200/80 text-base-content/70 border-0 font-medium">{category}</span>
			</div>

			<!-- Icon -->
			<div
				class="inline-flex size-11 items-center justify-center rounded-xl sm:size-14 lg:size-16 transition-transform duration-300 hover:scale-110 {iconBgColor}"
			>
				{#if iconColor}
					<div class={iconColor}>
						{@render icon()}
					</div>
				{:else}
					{@render icon()}
				{/if}
			</div>

			<!-- Title and Description -->
			<div class="flex-1 min-h-0">
				<h4 class="font-bold text-sm leading-tight mb-1.5 sm:text-base lg:text-lg line-clamp-2">{title}</h4>
				<p class="text-xs text-base-content/60 line-clamp-2 sm:text-sm">{description}</p>
			</div>
		</div>
	</div>
</a>
