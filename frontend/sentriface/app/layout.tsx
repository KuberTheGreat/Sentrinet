import { AuthProvider } from "@/context/AuthContext";
import "./globals.css";
import { ReactNode } from "react";

export const metadata = {
  title: "SentriNet",
  description: "AI-powered network sentinel",
};

export default function RootLayout({ children }: { children: ReactNode }) {
  return (
    <html lang="en">
      <body className="min-h-screen bg-gray-950 text-gray-100 font-sans">
        <AuthProvider>
            {children}
        </AuthProvider>
      </body>
    </html>
  );
}
