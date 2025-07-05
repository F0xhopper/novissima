"use client";

import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { toast } from "sonner";

export default function UploadPage() {
  const [isLoading, setIsLoading] = useState(false);
  const [formData, setFormData] = useState({
    textEnglish: "",
    textLatin: "",
    theme: "",
    imageSource: "",
    textSource: "",
    image: null as File | null,
  });
  const apiUrl = process.env.NEXT_PUBLIC_API_URL;

  const handleInputChange = (field: string, value: string) => {
    setFormData((prev) => ({ ...prev, [field]: value }));
  };

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (file) {
      setFormData((prev) => ({ ...prev, image: file }));
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsLoading(true);

    try {
      const formDataToSend = new FormData();
      formDataToSend.append("textEnglish", formData.textEnglish);
      formDataToSend.append("textLatin", formData.textLatin);
      formDataToSend.append("theme", formData.theme);
      formDataToSend.append("imageSource", formData.imageSource);
      formDataToSend.append("textSource", formData.textSource);

      if (formData.image) {
        formDataToSend.append("image", formData.image);
      }

      const response = await fetch(`${apiUrl}/content`, {
        method: "POST",
        body: formDataToSend,
      });

      if (response.ok) {
        const result = await response.json();
        toast.success(result.message);
        setFormData({
          textEnglish: "",
          textLatin: "",
          theme: "",
          imageSource: "",
          textSource: "",
          image: null,
        });
      } else {
        const error = await response.text();
        toast.error(error);
      }
    } catch (error) {
      toast.error("Failed to upload content");
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="min-h-screen p-8">
      <div className="max-w-2xl mx-auto">
        <Card>
          <CardHeader>
            <CardTitle>Upload Content</CardTitle>
            <CardDescription>
              Add new content to the Novissima database
            </CardDescription>
          </CardHeader>
          <CardContent>
            <form onSubmit={handleSubmit} className="space-y-6">
              <div className="space-y-2">
                <Label htmlFor="textEnglish">English Text *</Label>
                <Textarea
                  id="textEnglish"
                  placeholder="Enter the English text..."
                  value={formData.textEnglish}
                  onChange={(e) =>
                    handleInputChange("textEnglish", e.target.value)
                  }
                  required
                />
              </div>

              <div className="space-y-2">
                <Label htmlFor="textLatin">Latin Text *</Label>
                <Textarea
                  id="textLatin"
                  placeholder="Enter the Latin text..."
                  value={formData.textLatin}
                  onChange={(e) =>
                    handleInputChange("textLatin", e.target.value)
                  }
                  required
                />
              </div>

              <div className="space-y-2">
                <Label htmlFor="theme">Theme *</Label>
                <Select
                  value={formData.theme}
                  onValueChange={(value) => handleInputChange("theme", value)}
                >
                  <SelectTrigger>
                    <SelectValue placeholder="Select a theme" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="heaven">Heaven</SelectItem>
                    <SelectItem value="hell">Hell</SelectItem>
                    <SelectItem value="death">Death</SelectItem>
                    <SelectItem value="judgement">Judgement</SelectItem>
                  </SelectContent>
                </Select>
              </div>

              <div className="space-y-2">
                <Label htmlFor="imageSource">Image Source</Label>
                <Input
                  id="imageSource"
                  placeholder="Enter image source..."
                  value={formData.imageSource}
                  onChange={(e) =>
                    handleInputChange("imageSource", e.target.value)
                  }
                />
              </div>

              <div className="space-y-2">
                <Label htmlFor="textSource">Text Source</Label>
                <Input
                  id="textSource"
                  placeholder="Enter text source..."
                  value={formData.textSource}
                  onChange={(e) =>
                    handleInputChange("textSource", e.target.value)
                  }
                />
              </div>

              <div className="space-y-2">
                <Label htmlFor="image">Image</Label>
                <Input
                  id="image"
                  type="file"
                  accept="image/*"
                  onChange={handleFileChange}
                />
                <p>Maximum file size: 5MB. Supported formats: JPEG, PNG, GIF</p>
              </div>

              <Button type="submit" disabled={isLoading} className="w-full">
                {isLoading ? "Uploading..." : "Upload Content"}
              </Button>
            </form>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
