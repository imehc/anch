<script lang="ts">
	import type { Snippet } from 'svelte';

	interface SectionHeaderProps {
		title: string;
		icon: Snippet;
		iconBgColor: string;
		iconColor: string;
	}

	let { title, icon, iconBgColor, iconColor }: SectionHeaderProps = $props();
</script>

<div class="mb-4 flex items-center gap-2 sm:mb-6 sm:gap-3">
	<div
		class="flex size-10 items-center justify-center rounded-xl sm:size-12 sm:rounded-2xl {iconBgColor}"
	>
		<div class={iconColor}>
			{@render icon()}
		</div>
	</div>
	<h3 class="text-xl font-bold sm:text-2xl">{title}</h3>
</div>
