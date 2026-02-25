import { defineSchema, defineTable } from "convex/server";
import { v } from "convex/values";

export const schema = defineSchema({
  // the name of the table should match the name of the table in your convex dashboard
  todoList: defineTable({
    // here, the names of the columns in the table are defined, along with their types.
    body: v.string(),
    completed: v.boolean(),
  }),
});
