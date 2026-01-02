export interface Product {
  id: number;
  name: string;
  description: string;
  price: number;
  images: string;
  category: {
    id: number;
    name: string;
  };
  color_varian?: ColorVariant[];
}

export interface ColorVariant {
  id: number;
  name: string;
  color: string;
  images?: string;
  size_varian?: SizeVariant[];
}

export interface SizeVariant {
  id: number;
  size: string;
  stock: number;
}
