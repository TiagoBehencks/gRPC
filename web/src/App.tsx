import { useEffect, useState } from "react";
import type { Product } from "./proto/product";


function App() {
  const [products, setProducts] = useState<Product[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const [form, setForm] = useState({
    name: "",
    price: "",
    quantity: "",
  });

  const [editingId, setEditingId] = useState<string | null>(null);

  const fetchProducts = async () => {
    setLoading(true);
    setError(null);
    try {
      const res = await fetch(`/api/products`);
      const data = await res.json();
      setProducts(data);
    } catch (e) {
      setError(String(e));
    } finally {
      setLoading(false);
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError(null);

    try {
      const method = editingId ? "PUT" : "POST";
      const body = JSON.stringify({
        name: form.name,
        price: parseFloat(form.price) || 0,
        quantity: parseInt(form.quantity) || 0,
      });

      const url = editingId
        ? `/api/products/${editingId}`
        : `/api/products`;

      const res = await fetch(url, {
        method,
        headers: { "Content-Type": "application/json" },
        body,
      });

      if (!res.ok) throw new Error(await res.text());

      setEditingId(null);
      setForm({ name: "", price: "", quantity: "" });
      await fetchProducts();
    } catch (e) {
      setError(String(e));
    } finally {
      setLoading(false);
    }
  };

  const handleEdit = (product: Product) => {
    setEditingId(product.id);
    setForm({
      name: product.name,
      price: String(product.price),
      quantity: String(product.quantity),
    });
  };

  const handleDelete = async (id: string) => {
    setLoading(true);
    setError(null);
    try {
      const res = await fetch(`/api/products/${id}`, {
        method: "DELETE",
      });
      if (!res.ok) throw new Error(await res.text());
      await fetchProducts();
    } catch (e) {
      setError(String(e));
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchProducts();
  }, []);

  return (
    <div style={{ padding: "20px", maxWidth: "600px", margin: "0 auto" }}>
      <h1>Products CRUD</h1>

      <form
        onSubmit={handleSubmit}
        style={{ display: "flex", flexDirection: "column", gap: "10px", marginBottom: "20px" }}
      >
        <input
          type="text"
          placeholder="Name"
          value={form.name}
          onChange={(e) => setForm({ ...form, name: e.target.value })}
          required
          style={{ padding: "8px", fontSize: "16px" }}
        />
        <input
          type="number"
          placeholder="Price"
          value={form.price}
          onChange={(e) => setForm({ ...form, price: e.target.value })}
          required
          step="0.01"
          style={{ padding: "8px", fontSize: "16px" }}
        />
        <input
          type="number"
          placeholder="Quantity"
          value={form.quantity}
          onChange={(e) => setForm({ ...form, quantity: e.target.value })}
          required
          style={{ padding: "8px", fontSize: "16px" }}
        />
        <button
          type="submit"
          disabled={loading}
          style={{ padding: "10px", fontSize: "16px", cursor: "pointer" }}
        >
          {loading ? "Loading..." : editingId ? "Update" : "Add Product"}
        </button>
      </form>

      {error && <div style={{ color: "red", marginBottom: "10px" }}>Error: {error}</div>}

      <div>
        {products.map((p) => (
          <div
            key={p.id}
            style={{
              border: "1px solid #ccc",
              padding: "10px",
              marginBottom: "10px",
              display: "flex",
              justifyContent: "space-between",
              alignItems: "center",
            }}
          >
            <div>
              <strong>{p.name}</strong> - ${p.price} (qty: {p.quantity})
            </div>
            <div style={{ display: "flex", gap: "10px" }}>
              <button type="button" onClick={() => handleEdit(p)} style={{ cursor: "pointer" }}>
                Edit
              </button>
              <button
                type="button"
                onClick={() => handleDelete(p.id)}
                style={{ cursor: "pointer", color: "red" }}
              >
                Delete
              </button>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}

export default App;