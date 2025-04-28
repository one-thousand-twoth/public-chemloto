import polymers from '@/../polymers.json';

export type Polymers = typeof polymers;
export type Field = keyof Polymers;

type AlphaField = keyof Polymers['Альфа']
type BetaField = keyof Polymers['Бета']
type GammafField = keyof Polymers['Гамма']

export type CommonStructureNames = AlphaField | BetaField | GammafField
export type StructureNames<K extends Field> = keyof Polymers[K];

type Fields = {
	[key in Field]: {
		[polymerName in keyof Polymers[key]]: Polymer;
	};
};
export interface Polymer extends Array<Entry> {}
interface Entry {
	[element: string]: number;
}
/** Содержит в себе перечисление полей с вложеными структурами 
 * 
 * Десериализуется из файла polymers.json в корне проекта */ 
export const Polymers = polymers as Fields


