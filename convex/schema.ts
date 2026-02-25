import { defineSchema, defineTable } from "convex/server";
import { v } from "convex/values";

export default defineSchema({
  todos: defineTable({
    body: v.string(),
    completed: v.boolean(),
  }),
});
