'use client';

import { useState, useCallback } from "react";
import { Upload, FileImage, Download, X, Loader2 } from "lucide-react";

export default function Home() {
  const [files, setFiles] = useState<File[]>([]);
  const [previews, setPreviews] = useState<string[]>([]);
  const [isDragging, setIsDragging] = useState(false);
  const [isConverting, setIsConverting] = useState(false);

  const handleDrop = useCallback((e: React.DragEvent) => {
    e.preventDefault();
    setIsDragging(false);

    const droppedFiles = Array.from(e.dataTransfer.files).filter((file) =>
      file.type.startsWith("image/")
    );

    addFiles(droppedFiles);
  }, []);

  const addFiles = (newFiles: File[]) => {
    setFiles((prev) => [...prev, ...newFiles]);

    newFiles.forEach((file) => {
      const reader = new FileReader();
      reader.onload = (e) => {
        setPreviews((prev) => [...prev, e.target?.result as string]);
      };
      reader.readAsDataURL(file);
    });
  };

  const handleFileInput = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files) {
      addFiles(Array.from(e.target.files));
    }
  };

  const removeFile = (index: number) => {
    setFiles((prev) => prev.filter((_, i) => i !== index));
    setPreviews((prev) => prev.filter((_, i) => i !== index));
  };

  const convertToPDF = async () => {
    if (files.length === 0) return;

    setIsConverting(true);
    const formData = new FormData();
    files.forEach((file) => formData.append("files", file));

    try {
      const response = await fetch("http://localhost:8080/convert/image-pdf", {
        method: "POST",
        body: formData,
      });

      if (!response.ok) throw new Error("Conversion failed");

      const blob = await response.blob();
      const url = window.URL.createObjectURL(blob);
      const a = document.createElement("a");
      a.href = url;
      a.download = "converted.pdf";
      document.body.appendChild(a);
      a.click();
      window.URL.revokeObjectURL(url);
      document.body.removeChild(a);

      // Reset after successful conversion
      setFiles([]);
      setPreviews([]);
    } catch (error) {
      console.error("Error converting to PDF:", error);
      alert("Failed to convert images to PDF");
    } finally {
      setIsConverting(false);
    }
  };

  return (
    <div className="min-h-screen bg-white dark:bg-black">
      {/* Header */}
      <header className="border-b border-gray-200 dark:border-gray-800 bg-white dark:bg-black">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
          <div className="flex items-center gap-3">
            <div className="bg-black dark:bg-white p-2 rounded-lg">
              <FileImage className="w-6 h-6 text-white dark:text-black" />
            </div>
            <h1 className="text-2xl font-bold text-black dark:text-white">
              QuickConvert
            </h1>
          </div>
          <p className="mt-2 text-sm text-gray-600 dark:text-gray-400">
            Convert your images to PDF in seconds
          </p>
        </div>
      </header>

      {/* Main Content */}
      <main className="max-w-5xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
        {/* Upload Area */}
        <div
          onDrop={handleDrop}
          onDragOver={(e) => {
            e.preventDefault();
            setIsDragging(true);
          }}
          onDragLeave={() => setIsDragging(false)}
          className={`relative border-2 border-dashed rounded-xl p-12 text-center transition-all duration-200 ${
            isDragging
              ? "border-black dark:border-white bg-gray-50 dark:bg-gray-900"
              : "border-gray-300 dark:border-gray-700 bg-gray-50 dark:bg-gray-950 hover:border-gray-400 dark:hover:border-gray-600"
          }`}
        >
          <input
            type="file"
            multiple
            accept="image/*"
            onChange={handleFileInput}
            className="hidden"
            id="file-input"
          />
          <label htmlFor="file-input" className="cursor-pointer">
            <div className="flex flex-col items-center gap-4">
              <div className="bg-black dark:bg-white p-6 rounded-xl">
                <Upload className="w-12 h-12 text-white dark:text-black" />
              </div>
              <div>
                <h3 className="text-xl font-semibold text-black dark:text-white mb-2">
                  Drop your images here
                </h3>
                <p className="text-gray-600 dark:text-gray-400">
                  or click to browse from your device
                </p>
              </div>
              <div className="flex gap-2 text-xs text-gray-600 dark:text-gray-400">
                <span className="px-3 py-1 bg-white dark:bg-gray-900 border border-gray-200 dark:border-gray-800 rounded-full">
                  JPG
                </span>
                <span className="px-3 py-1 bg-white dark:bg-gray-900 border border-gray-200 dark:border-gray-800 rounded-full">
                  PNG
                </span>
                <span className="px-3 py-1 bg-white dark:bg-gray-900 border border-gray-200 dark:border-gray-800 rounded-full">
                  JPEG
                </span>
              </div>
            </div>
          </label>
        </div>

        {/* Preview Grid */}
        {previews.length > 0 && (
          <div className="mt-8">
            <div className="flex items-center justify-between mb-4">
              <h2 className="text-lg font-semibold text-black dark:text-white">
                Selected Images ({files.length})
              </h2>
              <button
                onClick={convertToPDF}
                disabled={isConverting}
                className="flex items-center gap-2 px-6 py-3 bg-black dark:bg-white text-white dark:text-black rounded-lg font-medium hover:bg-gray-800 dark:hover:bg-gray-200 transition-all duration-200 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                {isConverting ? (
                  <>
                    <Loader2 className="w-5 h-5 animate-spin" />
                    Converting...
                  </>
                ) : (
                  <>
                    <Download className="w-5 h-5" />
                    Convert to PDF
                  </>
                )}
              </button>
            </div>

            <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 gap-4">
              {previews.map((preview, index) => (
                <div
                  key={index}
                  className="relative group rounded-lg overflow-hidden bg-white dark:bg-gray-950 border border-gray-200 dark:border-gray-800 hover:border-black dark:hover:border-white transition-all duration-200"
                >
                  <img
                    src={preview}
                    alt={`Preview ${index + 1}`}
                    className="w-full h-48 object-cover"
                  />
                  <button
                    onClick={() => removeFile(index)}
                    className="absolute top-2 right-2 p-2 bg-black dark:bg-white text-white dark:text-black rounded-full opacity-0 group-hover:opacity-100 transition-opacity duration-200 hover:bg-gray-800 dark:hover:bg-gray-200"
                  >
                    <X className="w-4 h-4" />
                  </button>
                  <div className="p-2 text-xs text-gray-600 dark:text-gray-400 truncate">
                    {files[index].name}
                  </div>
                </div>
              ))}
            </div>
          </div>
        )}

        {/* Features */}
        <div className="mt-16 grid grid-cols-1 md:grid-cols-3 gap-6">
          {[
            {
              title: "Fast Conversion",
              desc: "Convert multiple images in seconds",
              icon: "⚡",
            },
            {
              title: "Secure",
              desc: "Your files are processed locally",
              icon: "🔒",
            },
            {
              title: "No Limits",
              desc: "Convert unlimited images to PDF",
              icon: "∞",
            },
          ].map((feature, i) => (
            <div
              key={i}
              className="p-6 rounded-lg bg-gray-50 dark:bg-gray-950 border border-gray-200 dark:border-gray-800"
            >
              <div className="text-3xl mb-3">{feature.icon}</div>
              <h3 className="font-semibold text-black dark:text-white mb-1">
                {feature.title}
              </h3>
              <p className="text-sm text-gray-600 dark:text-gray-400">
                {feature.desc}
              </p>
            </div>
          ))}
        </div>
      </main>
    </div>
  );
}
