# Descriptions: Good vs Bad

The `description` is the **only** signal the agent uses to decide whether to trigger a skill. A vague description → skill never fires. A specific description with concrete keywords → skill fires at the right moment.

## Rules

1. **Third person.** "Authors skills..." not "You author skills...".
2. **What + when.** Both must be present.
3. **Concrete keywords.** Names of tools, formats, file types, domains the user is likely to mention.
4. **≤ 2 sentences.**

## Side-by-side examples

### Too vague
> Helps with documents.

Never triggers reliably. "Documents" is too generic.

### Better
> Extracts text and tables from PDF files, fills forms, merges documents. Use when working with PDF files or when the user mentions PDFs, forms, or document extraction.

Concrete keywords ("PDF", "forms", "extraction"); what + when both present.

---

### Too vague
> Cape Town stuff.

### Better
> Reference material on Cape Town's City Bowl and Atlantic Seaboard neighbourhoods — real landmarks, history, culture, typical activity — used to inform NPC dialogue, quest design, building names, and atmosphere. Use when writing dialogue or quests, naming buildings, or choosing which POIs are interactive in a given area.

---

### Wrong person (second)
> Use this skill when you want to add a Phaser scene.

### Right (third)
> Phaser 3 patterns for this game: scene lifecycle, Tiled map loading, player movement with arrow/WASD keys, camera follow, layer-based collision, NPC interaction, and dialogue. Use when wiring new game code, adding a scene, debugging movement or collision, or integrating a new interactive object type.

---

### Missing "when"
> Defines the game's pixel-art conventions.

### Complete
> Defines the game's pixel-art conventions — 32×32 tile grid, 3/4 top-down perspective, palette, walk-cycle frame counts, building-facade construction — and guides authoring or commissioning new sprites. Use when adding a character, building, prop, or tileset art, or when reviewing art for consistency.

---

### Missing "what" (just categories)
> For game evaluation tasks.

### Complete
> Evaluation scenarios that exercise the game end-to-end — walk-route coverage, collision correctness, POI trigger firing, visual consistency, dialogue flow. Use when adding a feature, before merging map or engine changes, or when building a new skill.

## The trigger test

Ask: "If the agent had 100 skills listed by description alone, and a user said `<realistic prompt>`, would it pick this one over the others?"

If the answer isn't obvious, make the description more specific — usually by naming the tool, format, or domain term the user is most likely to say.
