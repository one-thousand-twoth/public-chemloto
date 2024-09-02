import polymers from '@/../../polymers.json' 
/** Содержит в себе перечисление полей с вложеными структурами 
 * 
 * Десериализуется из файла polymers.json в корне проекта */ 
interface Fields {
	[key: string]: {
		[polymerName: string]: Polymer;
	};
}
interface Polymer extends Array<Entry> {}
interface Entry {
	[element: string]: number;
}

export const Polymers = polymers as Fields


