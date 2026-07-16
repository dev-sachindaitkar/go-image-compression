import { Inter } from "next/font/google"; // ◄ Import a premium, high-readability font
import "./globals.css";
import { Metadata } from "next";

// Configure font weights for clean title headers and micro data text
const inter = Inter({
  subsets: ["latin"],
  weight: ["400", "500", "600", "700", "800"],
  variable: "--font-sans",
});

export const metadata: Metadata = {
  title: "Image Compressor",
  description: "High-performance parallel image optimization powered by Go",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en" className={`${inter.variable}`}>
      <body className="font-sans"> {children}</body>
    </html>
  );
}
