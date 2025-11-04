// frontend/sentriface/app/login/page.tsx
"use client";
import { useForm } from "react-hook-form";
import API from "@/lib/api";
import { useAuth } from "@/context/AuthContext";
import { useRouter } from "next/navigation";

export default function LoginPage() {
  const { register, handleSubmit } = useForm();
  const { setToken } = useAuth();
  const router = useRouter();

  const onSubmit = async (data: any) => {
    try {
      const res = await API.post("/login", data);
      setToken(res.data.token);
      router.push("/dashboard");
    } catch (err) {
      alert("Login failed");
    }
  };

  return (
    <div className="flex flex-col items-center justify-center h-screen">
      <h1 className="text-2xl mb-4 font-semibold">Login</h1>
      <form onSubmit={handleSubmit(onSubmit)} className="flex flex-col gap-3">
        <input {...register("email")} placeholder="Email" className="border p-2 rounded" />
        <input {...register("password")} placeholder="Password" type="password" className="border p-2 rounded" />
        <button type="submit" className="bg-blue-600 text-white px-4 py-2 rounded">Login</button>
      </form>
    </div>
  );
}
