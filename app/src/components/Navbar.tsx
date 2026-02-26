import { Box, Flex, Button, Text, Container } from "@chakra-ui/react";
import { IoMoon } from "react-icons/io5";

export default function Navbar() {
	return (
		<Container maxW={"900px"}>
			<Box bg= {( "gray.700")} px={4} my={4} borderRadius={"5"}>
				<Flex h={16} alignItems={"center"} justifyContent={"space-between"}>
					{/* LEFT SIDE */}
					<Flex
						justifyContent={"center"}
						alignItems={"center"}
						gap={3}
						display={{ base: "none", sm: "flex" }}
					>
						{/*<img src='/react.png' alt='logo' width={50} height={50} />*/}
						{/*<Text fontSize={"40"}>+</Text>*/}
						{/*<img src='/go.png' alt='logo' width={40} height={40} />*/}
						{/*<Text fontSize={"40"}>=</Text>*/}
						{/*<img src='/explode.png' alt='logo' width={50} height={50} />*/}
					</Flex>

					{/* RIGHT SIDE */}
					<Flex alignItems={"center"} gap={3}>
						<Text fontSize={"lg"} fontWeight={500}>
							Daily Tasks
						</Text>
						<Button>
							<IoMoon />
						</Button>
					</Flex>
				</Flex>
			</Box>
		</Container>
	);
}
