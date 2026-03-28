# Design System Documentation: The Architectural Blueprint

## 1. Overview & Creative North Star
### The Creative North Star: "Precision Editorial"
This design system rejects the "generic dashboard" aesthetic in favor of a high-density, editorial approach. We are building for power users who require speed and clarity, but we refuse to sacrifice soul. Our goal is to balance the developer-friendly logic of Mantine UI with the sophisticated, layered depth of a premium digital workspace.

We achieve this through **Precision Editorial**—a philosophy where information density is handled not through lines and boxes, but through intentional "Tonal Layering." By utilizing a tight typographic scale and a sophisticated indigo-based dark palette, we create an environment that feels like a high-end IDE crossed with a modern architectural journal.

---

## 2. Colors & Surface Philosophy
The palette is rooted in deep indigo and cool slates, designed to reduce eye strain while maintaining a vibrant "primary" energy.

### The "No-Line" Rule
**Borders are a design failure of the past.** To maintain a premium feel, designers are prohibited from using `1px solid` borders for sectioning. Structural separation must be achieved through:
1.  **Background Color Shifts:** Moving from `surface` to `surface_container_low`.
2.  **Negative Space:** Using the Spacing Scale (specifically `spacing.4` to `spacing.8`).

### Surface Hierarchy & Nesting
Treat the UI as a series of physical layers. Each "inner" container should move up or down the tier to define its importance:
-   **Base Layer:** `surface` (#121316)
-   **Structural Blocks:** `surface_container_low` (#1a1b1e)
-   **Interactive Cards:** `surface_container` (#1f1f23)
-   **Floating Popovers/Modals:** `surface_container_high` (#292a2d)

### The Glass & Gradient Rule
To prevent the dark mode from feeling "flat," use **Glassmorphism** for floating elements (Top Nav, Context Menus).
-   **Recipe:** `surface_container` at 80% opacity + `backdrop-blur: 12px`.
-   **Signature Accent:** Use a subtle linear gradient for Primary CTAs: `brand` (#ae90da) to `brand_deep` (#654c96) at a 135-degree angle. This adds a "lithographic" depth that flat hex codes lack.

---

## 3. Typography
We utilize **Inter** for its mathematical precision and neutral character, allowing the content to lead.

*   **Display (lg/md):** Reserved for hero data points or section starters. Use `display-md` (2.75rem) with `-0.02em` letter-spacing to feel "tight" and authoritative.
*   **Headlines & Titles:** Use `headline-sm` (1.5rem) for major module headers. Pair with `primary_fixed` color to pull the eye.
*   **Body & Labels:** This is where density lives. `body-md` (0.875rem) is our workhorse. For secondary metadata, use `label-md` (0.75rem) in `on_surface_variant` to create clear visual hierarchy without changing font sizes.

---

## 4. Elevation & Depth
In this system, depth is **Tonal**, not structural.

### The Layering Principle
Avoid "Drop Shadows" on standard cards. Instead, place a `surface_container_lowest` (#0d0e11) element inside a `surface_container` (#1f1f23) area to create a "recessed" effect. This mimics a physical milled surface.

### Ambient Shadows
When an element must float (e.g., a Modal), use a **Tinted Ambient Shadow**:
-   **Blur:** 20px to 40px.
-   **Spread:** -5px.
-   **Color:** `rgba(0, 0, 0, 0.4)` mixed with 4% of `primary`. This prevents the shadow from looking "dirty" on the deep indigo background.

### The "Ghost Border" Fallback
If a border is required for accessibility (e.g., Input fields), use a **Ghost Border**: `outline_variant` (#404752) at **20% opacity**. It should be felt, not seen.

---

## 5. Components

### Buttons
-   **Primary:** Gradient fill (`primary` to `primary_container`), `radius.md` (0.75rem). Text: `on_primary`.
-   **Secondary:** Ghost style. No background, `outline_variant` @ 20% border. On hover, shift background to `surface_container_high`.
-   **Tertiary:** Text-only using `primary` color. No box.

### Input Fields
-   **Surface:** `surface_container_low`.
-   **Interaction:** On focus, the "Ghost Border" becomes `primary` at 100% opacity, and a subtle `primary` outer glow (4px blur, 10% opacity) appears.
-   **Density:** Use `spacing.3` (0.6rem) for vertical padding to keep the UI compact.

### Cards & Lists
-   **Forbid Dividers:** Never use a horizontal line to separate list items. Use a 1px gap showing the `surface_container_lowest` background, or simply use `spacing.2` of vertical rhythm.
-   **Roundedness:** All cards must use `radius.lg` (1rem). Nested items (like inner buttons) must use `radius.md` (0.75rem) to create proper visual nesting (the "inner radius < outer radius" rule).

### Additional Component: The "Action Bar"
A floating, glassmorphic bar (bottom-center) for bulk actions.
-   **Style:** `surface_container_highest` @ 70% opacity, `backdrop-blur: 20px`, `radius.full`.
-   **Purpose:** Keeps the main workspace clean by hiding secondary actions until items are selected.

---

## 6. Do's and Don'ts

### Do
*   **Do** use `surface_container` variations to group related data instead of drawing boxes.
*   **Do** lean into high-contrast typography (e.g., `primary` headers with `on_surface_variant` body text).
*   **Do** use Mantine’s `lg` (1rem) rounding for outer containers to soften the "technical" feel of the indigo palette.

### Don't
*   **Don't** use 100% black (#000000) or pure grey. Always use the indigo-tinted `surface` tokens.
*   **Don't** use `spacing.16` or larger unless it's a landing page. For "Developer-Friendly" density, keep gaps between `spacing.4` and `spacing.8`.
*   **Don't** use standard "drop shadows" with 0 blur. Shadows must be expansive and atmospheric.
