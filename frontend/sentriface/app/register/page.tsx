"use client";

import { useForm } from "react-hook-form";
import { useRouter } from "next/navigation";
import { useAuth } from "@/context/AuthContext";
import { useState } from "react";

type RegisterForm = {
  username: string;
  password: string;
};

export default function RegisterPage() {
  const { register, handleSubmit } = useForm<RegisterForm>();
  const { login } = useAuth();
  const router = useRouter();
  const [loading, setLoading] = useState(false);

  const onSubmit = async (data: RegisterForm) => {
    setLoading(true);
    try {
      const res = await fetch("https://sentrinet.onrender.com/register", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(data),
      });

      const result = await res.json();

      if (!res.ok) throw new Error(result?.message || "Registration failed");

      if (!result.token) throw new Error("No token in response");

      login(result.token);
      router.push("/dashboard");
    } catch (err: any) {
      alert(err.message || "Error during registration");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen from-blue-500 via-indigo-500 to-purple-500 flex items-center justify-center p-6">
      <div className="bg-gray-800 backdrop-blur-lg rounded-2xl shadow-2xl p-8 w-full max-w-md transition-all">
        <h1 className="text-3xl font-semibold text-center mb-6 text-gray-800">
          Create Your Account
        </h1>

        <form onSubmit={handleSubmit(onSubmit)} className="flex flex-col gap-5">
          <input
            {...register("username")}
            placeholder="Username"
            className="p-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500 transition-all"
            required
          />

          <input
            {...register("password")}
            placeholder="Password"
            type="password"
            className="p-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500 transition-all"
            required
          />

          <button
            type="submit"
            disabled={loading}
            className={`py-3 text-white rounded-lg font-medium transition-all ${
              loading
                ? "bg-indigo-300 cursor-not-allowed"
                : "bg-indigo-600 hover:bg-indigo-700"
            }`}
          >
            {loading ? "Registering..." : "Sign Up"}
          </button>
        </form>

        <p className="text-sm text-gray-600 text-center mt-5">
          Already registered?{" "}
          <span
            onClick={() => router.push("/login")}
            className="text-indigo-600 hover:underline cursor-pointer"
          >
            Log in
          </span>
        </p>
      </div>
    </div>
  );
}
