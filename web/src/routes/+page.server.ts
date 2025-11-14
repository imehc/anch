import type { PageServerLoad } from './$types';
import * as m from '$lib/paraglide/messages';

export type IconName = 'Play' | 'Sliders' | 'Palette';

export interface ItemData {
	title: string;
	href: string;
	description: string;
	iconType: IconName | 'number';
	iconContent?: string;
	iconBgColor: string;
	iconColor?: string;
	category: string;
}

export interface SectionData {
	title: string;
	iconType: 'Sliders' | 'Play';
	iconBgColor: string;
	iconColor: string;
	items: ItemData[];
}

function getAllSections(): SectionData[] {
	return [
		{
			title: m.section_tools(),
			iconType: 'Sliders',
			iconBgColor: 'bg-blue-100',
			iconColor: 'text-blue-500',
			items: [
				{
					title: m.tool_color_title(),
					href: '/tools/color',
					description: m.tool_color_desc(),
					iconType: 'Palette',
					iconBgColor: 'bg-yellow-100',
					iconColor: 'text-yellow-600',
					category: m.category_design()
				}
			]
		},
		// {
		// 	title: m.section_games(),
		// 	iconType: 'Play',
		// 	iconBgColor: 'bg-purple-100',
		// 	iconColor: 'text-purple-500',
		// 	items: [
		// 		{
		// 			title: m.game_2048_title(),
		// 			href: '/games/2048',
		// 			description: m.game_2048_desc(),
		// 			iconType: 'number',
		// 			iconContent: '2',
		// 			iconBgColor: 'bg-orange-100',
		// 			iconColor: 'text-orange-600',
		// 			category: m.category_puzzle()
		// 		},
		// 	]
		// }
	];
}

export const load: PageServerLoad = async ({ url }) => {
	const query = url.searchParams.get('q')?.toLowerCase().trim() || '';

	const allSections = getAllSections();

	if (!query) {
		return { sections: allSections, query };
	}

	// Filter sections based on search query
	const filteredSections = allSections
		.map((section) => ({
			...section,
			items: section.items.filter(
				(item) =>
					item.title.toLowerCase().includes(query) ||
					item.description.toLowerCase().includes(query) ||
					item.category.toLowerCase().includes(query)
			)
		}))
		.filter((section) => section.items.length > 0);

	return { sections: filteredSections, query };
};
