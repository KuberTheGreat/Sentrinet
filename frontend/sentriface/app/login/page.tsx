"use client";

import { useForm } from "react-hook-form";
import { useAuth } from "@/context/AuthContext";
import { useRouter } from "next/navigation";

export default function LoginPage() {
  const { register, handleSubmit } = useForm();
  const { login } = useAuth();
  const router = useRouter();

  const onSubmit = async (data: any) => {
    try {
      const res = await fetch("https://sentrinet.onrender.com/login", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          username: data.username,
          password: data.password,
        }),
      });

      if (!res.ok) {
        alert("Invalid credentials");
        return;
      }

      const response = await res.json();
      if (response.token) {
        login(response.token);
        router.push("/dashboard");
      } else {
        alert("Login failed: No token received");
      }
    } catch (err) {
      console.error("Login error:", err);
      alert("Something went wrong during login");
    }
  };

  return (
    <div className="flex flex-col items-center justify-center h-screen bg-gray-950 text-white">
      <h1 className="text-3xl mb-6 font-semibold">Welcome Back 👋</h1>

      <form
        onSubmit={handleSubmit(onSubmit)}
        className="flex flex-col gap-4 bg-gray-900 p-8 rounded-xl shadow-lg w-80"
      >
        <input
          {...register("username")}
          placeholder="Username"
          className="border border-gray-700 bg-gray-800 p-3 rounded focus:outline-none focus:ring focus:ring-blue-500"
        />
        <input
          {...register("password")}
          placeholder="Password"
          type="password"
          className="border border-gray-700 bg-gray-800 p-3 rounded focus:outline-none focus:ring focus:ring-blue-500"
        />

        <button
          type="submit"
          className="bg-blue-600 hover:bg-blue-700 text-white font-semibold py-2 rounded transition"
        >
          Log In
        </button>
      </form>

      <p className="mt-4 text-gray-400 text-sm">
        Don’t have an account?{" "}
        <a href="/register" className="text-blue-400 hover:text-blue-300 underline">
          Register here
        </a>
      </p>
    </div>
  );
}
