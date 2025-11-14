<script lang="ts">
	let { class: className = 'size-6' }: { class?: string } = $props();
</script>

<svg class={className} fill="none" stroke="currentColor" viewBox="0 0 24 24">
	<path
		stroke-linecap="round"
		stroke-linejoin="round"
		stroke-width="2"
		d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
	/>
</svg>
