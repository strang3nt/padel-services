export function genderToString(g: number): string {
  switch (g) {
    case 0:
      return "Male";
    case 1:
      return "Female";
    case 2:
      return "Mixed";
    default:
      throw Error("Team gender not recognized");
  }
}
