import { z } from "zod";

export const UserSchema = z.object({
  ID: z.string(),
  Email: z.string(),
  Name: z.string(),
  SessionID: z.string(),
});

export type UserSchemaType = z.infer<typeof UserSchema>;
