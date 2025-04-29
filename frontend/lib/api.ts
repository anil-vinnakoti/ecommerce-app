// lib/api.ts
export async function fetchProducts() {
  const res = await fetch("http://localhost:8080/products", {
    next: { revalidate: 60 } // optional: ISR
  });

  if (!res.ok) {
    throw new Error("Failed to fetch products");
  }

  return res.json();
}
