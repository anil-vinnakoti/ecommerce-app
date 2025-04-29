// app/products/page.tsx
import { fetchProducts } from "@/lib/api";

export default async function ProductsPage() {
  const products = await fetchProducts();

  return (
    <div className="grid grid-cols-1 md:grid-cols-3 gap-4 p-4">
      {products.map((product: any) => (
        <div key={product.ID} className="border p-4 rounded shadow">
          <img
            src={product.image_url}
            alt={product.name}
            className="w-full h-40 object-cover rounded"
          />
          <h2 className="text-xl font-semibold">{product.name}</h2>
          <p className="text-gray-700">{product.description}</p>
          <p className="text-green-600 font-bold">${product.price}</p>
        </div>
      ))}
    </div>
  );
}
